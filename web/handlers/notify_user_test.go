package handlers_test

import (
    "bufio"
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net"
    "net/http"
    "net/http/httptest"
    "os"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/web/handlers"
    "github.com/dgrijalva/jwt-go"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("NotifyUser", func() {
    var fakeUAA *httptest.Server
    var buffer *bytes.Buffer
    var logger *log.Logger
    var writer *httptest.ResponseRecorder
    var request *http.Request
    var handler handlers.NotifyUser
    var envelope Envelope
    var smtpServer *net.TCPListener
    var token string

    BeforeEach(func() {
        tokenHeader := jwt.EncodeSegment([]byte(`{"alg":"RS256"}`))
        tokenBody := jwt.EncodeSegment([]byte(`{"jti":"c5f6a266-5cf0-4ae2-9647-2615e7d28fa1","client_id":"mister-client","cid":"mister-client","exp":1404281214}`))
        token = tokenHeader + "." + tokenBody

        fakeUAA = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
            if req.URL.Path == "/Users/user-123" &&
                req.Method == "GET" &&
                req.Header.Get("Authorization") == "Bearer "+token {

                response := `{
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

                w.WriteHeader(http.StatusOK)
                w.Write([]byte(response))
            } else {
                w.WriteHeader(http.StatusNotFound)
            }
        }))

        os.Setenv("UAA_HOST", fakeUAA.URL)
        os.Setenv("UAA_CLIENT_ID", "notifications")
        os.Setenv("UAA_CLIENT_SECRET", "secret")

        addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
        if err != nil {
            panic(err)
        }

        smtpServer, err = net.ListenTCP("tcp", addr)
        if err != nil {
            panic(err)
        }

        parts := strings.SplitN(smtpServer.Addr().String(), ":", 2)
        host := parts[0]
        port := parts[1]

        go func() {
            conn, err := smtpServer.AcceptTCP()
            if err != nil {
                smtpServer.Close()
                return
            }
            input := bufio.NewReader(conn)
            output := bufio.NewWriter(conn)

            output.WriteString(fmt.Sprintf("220 %s\r\n", host))
            output.Flush()

            for {
                request, _ := input.ReadString('\n')
                response, exit := envelope.Respond(request)
                if response != "" {
                    output.WriteString(response + "\r\n")
                    output.Flush()
                }
                if exit {
                    break
                }
            }

            conn.Close()
        }()

        os.Setenv("SMTP_USER", "smtp-user")
        os.Setenv("SMTP_PASS", "smtp-pass")
        os.Setenv("SMTP_HOST", host)
        os.Setenv("SMTP_PORT", port)

        os.Setenv("SENDER", "test-user@example.com")

        envelope = Envelope{}
        buffer = bytes.NewBuffer([]byte{})
        logger = log.New(buffer, "", 0)
        writer = httptest.NewRecorder()
        handler = handlers.NewNotifyUser(logger)
    })

    AfterEach(func() {
        fakeUAA.Close()
        smtpServer.Close()
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

            Expect(envelope).To(Equal(Envelope{
                AuthUser: "smtp-user",
                AuthPass: "smtp-pass",
                From:     "<test-user@example.com>",
                To:       "<fake-user@example.com>",
                Data: []string{
                    "From: test-user@example.com",
                    "To: fake-user@example.com",
                    "Subject: CF Notification: Reset your password",
                    "X-CF-Client-ID: mister-client",
                    "Body:",
                    "",
                    `The following "Password reminder" notification was sent to you directly by the "Login system" component of Cloud Foundry:`,
                    "",
                    "Please reset your password by clicking on this link...",
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
