package handlers_test

import (
    "bytes"
    "encoding/json"
    "errors"
    "io/ioutil"
    "log"
    "net/http"
    "net/http/httptest"
    "os"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/mail"
    "github.com/cloudfoundry-incubator/notifications/web/handlers"
    "github.com/dgrijalva/jwt-go"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

const responseFor123 = `{
    "id": "user-123",
    "meta": {
        "version": 6,
        "created": "2014-05-22T22:36:36.941Z",
        "lastModified": "2014-06-25T23:10:03.845Z"
    },
    "userName": "admin",
    "name": {
        "familyName": "Admin",
        "givenName": "Mister"
    },
    "emails": [
        {
            "value": "fake-user@example.com"
        }
    ],
    "groups": [
        {
            "value": "e7f74565-4c7e-44ba-b068-b16072cbf08f",
            "display": "clients.read",
            "type": "DIRECT"
        }
    ],
    "approvals": [],
    "active": true,
    "verified": false,
    "schemas": [
        "urn:scim:schemas:core:1.0"
    ]
}`

const responseFor456 = `{
    "id": "user-456",
    "meta": {
        "version": 6,
        "created": "2014-05-22T22:36:36.941Z",
        "lastModified": "2014-06-25T23:10:03.845Z"
    },
    "userName": "bounce",
    "name": {
        "familyName": "Bounce",
        "givenName": "Mister"
    },
    "emails": [
        {
            "value": "bounce@example.com"
        }
    ],
    "groups": [
        {
            "value": "e7f74565-4c7e-44ba-b068-b16072cbf08f",
            "display": "clients.read",
            "type": "DIRECT"
        }
    ],
    "approvals": [],
    "active": true,
    "verified": false,
    "schemas": [
        "urn:scim:schemas:core:1.0"
    ]
}`

type FakeMailClient struct {
    messages    []mail.Message
    errorOnSend bool
}

func (fake *FakeMailClient) Connect() error {
    return nil
}

func (fake *FakeMailClient) Send(msg mail.Message) error {
    if fake.errorOnSend {
        return errors.New("BOOM!")
    }

    fake.messages = append(fake.messages, msg)
    return nil
}

var _ = Describe("NotifyUser", func() {
    var fakeUAA *httptest.Server
    var buffer *bytes.Buffer
    var logger *log.Logger
    var writer *httptest.ResponseRecorder
    var request *http.Request
    var handler handlers.NotifyUser
    var token string
    var client FakeMailClient

    BeforeEach(func() {
        tokenHeader := jwt.EncodeSegment([]byte(`{"alg":"RS256"}`))
        tokenBody := jwt.EncodeSegment([]byte(`{"jti":"c5f6a266-5cf0-4ae2-9647-2615e7d28fa1","client_id":"mister-client","cid":"mister-client","exp":1404281214}`))
        token = tokenHeader + "." + tokenBody

        fakeUAA = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
            if req.Method == "GET" && req.Header.Get("Authorization") == "Bearer "+token {
                switch req.URL.Path {
                case "/Users/user-123":
                    w.WriteHeader(http.StatusOK)
                    w.Write([]byte(responseFor123))
                case "/Users/user-456":
                    w.WriteHeader(http.StatusOK)
                    w.Write([]byte(responseFor456))
                default:
                    w.WriteHeader(http.StatusNotFound)
                }
            } else {
                w.WriteHeader(http.StatusNotFound)
            }
        }))

        os.Setenv("UAA_HOST", fakeUAA.URL)
        os.Setenv("UAA_CLIENT_ID", "notifications")
        os.Setenv("UAA_CLIENT_SECRET", "secret")
        os.Setenv("SENDER", "test-user@example.com")

        buffer = bytes.NewBuffer([]byte{})
        logger = log.New(buffer, "", 0)
        writer = httptest.NewRecorder()

        client = FakeMailClient{}
        handler = handlers.NewNotifyUser(logger, &client)
    })

    AfterEach(func() {
        fakeUAA.Close()
    })

    Context("when the request is valid", func() {
        BeforeEach(func() {
            requestBody, err := json.Marshal(map[string]string{
                "kind":               "forgot_password",
                "kind_description":   "Password reminder",
                "source_description": "Login system",
                "subject":            "Reset your password",
                "text":               "Please reset your password by clicking on this link...",
            })
            if err != nil {
                panic(err)
            }

            request, err = http.NewRequest("POST", "/users/user-123", bytes.NewReader(requestBody))
            if err != nil {
                panic(err)
            }
            request.Header.Set("Authorization", "Bearer "+token)
        })

        It("logs the email address of the recipient", func() {
            handler.ServeHTTP(writer, request)

            Expect(buffer.String()).To(ContainSubstring("Sending email to fake-user@example.com"))
        })

        It("logs the message envelope", func() {
            handler.ServeHTTP(writer, request)

            data := []string{
                "From: test-user@example.com",
                "To: fake-user@example.com",
                "Subject: CF Notification: Reset your password",
                "Body:",
                `The following "Password reminder" notification was sent to you directly by the "Login system" component of Cloud Foundry:`,
                "Please reset your password by clicking on this link...",
            }
            results := strings.Split(buffer.String(), "\n")
            for _, item := range data {
                Expect(results).To(ContainElement(item))
            }
        })

        It("talks to the SMTP server, sending the email", func() {
            handler.ServeHTTP(writer, request)

            Expect(len(client.messages)).To(Equal(1))

            msg := client.messages[0]
            Expect(msg).To(Equal(mail.Message{
                From:    "test-user@example.com",
                To:      "fake-user@example.com",
                Subject: "CF Notification: Reset your password",
                Body: `The following "Password reminder" notification was sent to you directly by the "Login system" component of Cloud Foundry:

Please reset your password by clicking on this link...`,
                Headers: []string{
                    "X-CF-Client-ID: mister-client",
                },
            }))
        })

        It("returns a status response for the sent mail", func() {
            handler.ServeHTTP(writer, request)

            Expect(writer.Code).To(Equal(http.StatusOK))
            parsed := []map[string]string{}
            err := json.Unmarshal(writer.Body.Bytes(), &parsed)
            if err != nil {
                panic(err)
            }

            Expect(parsed).To(Equal([]map[string]string{
                map[string]string{
                    "status": "delivered",
                },
            }))
        })

        Context("when the SMTP server fails to deliver the mail", func() {
            BeforeEach(func() {
                request.URL.Path = "/users/user-456"
            })

            It("returns a status indicating that delivery failed", func() {
                client.errorOnSend = true
                handler.ServeHTTP(writer, request)

                Expect(writer.Code).To(Equal(http.StatusOK))
                parsed := []map[string]string{}
                err := json.Unmarshal(writer.Body.Bytes(), &parsed)
                if err != nil {
                    panic(err)
                }

                Expect(parsed).To(Equal([]map[string]string{
                    map[string]string{
                        "status": "failed",
                    },
                }))
            })
        })

        PContext("when the SMTP server cannot be reached", func() {
            It("returns a status indicating that the server is unavailable", func() {
                handler.ServeHTTP(writer, request)

                Expect(writer.Code).To(Equal(http.StatusOK))
                parsed := []map[string]string{}
                err := json.Unmarshal(writer.Body.Bytes(), &parsed)
                if err != nil {
                    panic(err)
                }

                Expect(parsed).To(Equal([]map[string]string{
                    map[string]string{
                        "status": "unavailable",
                    },
                }))
            })
        })
    })

    Context("when the request is invalid", func() {
        BeforeEach(func() {
            requestBody, err := json.Marshal(map[string]string{
                "kind_description":   "Password reminder",
                "source_description": "Login system",
                "subject":            "Reset your password",
            })
            if err != nil {
                panic(err)
            }

            request, err = http.NewRequest("POST", "/users/user-123", bytes.NewReader(requestBody))
            if err != nil {
                panic(err)
            }
            request.Header.Set("Authorization", "Bearer "+token)
        })

        It("returns an error message", func() {
            handler.ServeHTTP(writer, request)

            parsed := map[string][]string{}
            err := json.Unmarshal(writer.Body.Bytes(), &parsed)
            if err != nil {
                panic(err)
            }

            Expect(parsed["errors"]).To(ContainElement(`"kind" is a required field`))
            Expect(parsed["errors"]).To(ContainElement(`"text" is a required field`))
        })

        Context("when the request body is missing", func() {
            BeforeEach(func() {
                request.Body = ioutil.NopCloser(bytes.NewReader([]byte{}))
            })

            It("returns an error message", func() {
                handler.ServeHTTP(writer, request)

                parsed := map[string][]string{}
                err := json.Unmarshal(writer.Body.Bytes(), &parsed)
                if err != nil {
                    panic(err)
                }

                Expect(parsed["errors"]).To(ContainElement(`"kind" is a required field`))
                Expect(parsed["errors"]).To(ContainElement(`"text" is a required field`))
            })
        })
    })
})
