package handlers_test

import (
    "bytes"
    "encoding/json"
    "errors"
    "io/ioutil"
    "log"
    "net/http"
    "net/http/httptest"
    "net/url"
    "os"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/mail"
    "github.com/cloudfoundry-incubator/notifications/web/handlers"
    "github.com/nu7hatch/gouuid"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("NotifyUser", func() {
    var handler handlers.NotifyUser
    var buffer *bytes.Buffer
    var logger *log.Logger
    var writer *httptest.ResponseRecorder
    var request *http.Request
    var token string
    var mailClient FakeMailClient
    var uaaClient FakeUAAClient

    BeforeEach(func() {
        tokenHeader := map[string]interface{}{
            "alg": "RS256",
        }
        tokenClaims := map[string]interface{}{
            "jti":       "c5f6a266-5cf0-4ae2-9647-2615e7d28fa1",
            "client_id": "mister-client",
            "cid":       "mister-client",
            "exp":       3404281214,
            "scope":     []string{"notifications.write"},
        }
        token = BuildToken(tokenHeader, tokenClaims)

        os.Setenv("SENDER", "test-user@example.com")

        buffer = bytes.NewBuffer([]byte{})
        logger = log.New(buffer, "", 0)
        writer = httptest.NewRecorder()

        mailClient = FakeMailClient{}
        uaaClient = FakeUAAClient{
            UsersByID: map[string]uaa.User{
                "user-123": uaa.User{
                    ID:       "user-123",
                    Username: "admin",
                    Name: uaa.Name{
                        FamilyName: "Admin",
                        GivenName:  "Mister",
                    },
                    Emails:   []string{"fake-user@example.com"},
                    Active:   true,
                    Verified: false,
                },
                "user-456": uaa.User{
                    ID:       "user-456",
                    Username: "bounce",
                    Name: uaa.Name{
                        FamilyName: "Bounce",
                        GivenName:  "Mister",
                    },
                    Emails:   []string{"bounce@example.com"},
                    Active:   true,
                    Verified: false,
                },
            },
        }

        guidGenerator := handlers.GUIDGenerationFunc(func() (*uuid.UUID, error) {
            guid := uuid.UUID([16]byte{0xDE, 0xAD, 0xBE, 0xEF, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55})
            return &guid, nil
        })

        handler = handlers.NewNotifyUser(logger, &mailClient, &uaaClient, guidGenerator)
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

            Expect(len(mailClient.messages)).To(Equal(1))

            msg := mailClient.messages[0]
            Expect(msg).To(Equal(mail.Message{
                From:    "test-user@example.com",
                To:      "fake-user@example.com",
                Subject: "CF Notification: Reset your password",
                Body: `The following "Password reminder" notification was sent to you directly by the "Login system" component of Cloud Foundry:

Please reset your password by clicking on this link...`,
                Headers: []string{
                    "X-CF-Client-ID: mister-client",
                    "X-CF-Notification-ID: deadbeef-aabb-ccdd-eeff-001122334455",
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
                mailClient.errorOnSend = true
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

        Context("when the SMTP server cannot be reached", func() {
            It("returns a status indicating that the server is unavailable", func() {
                mailClient.errorOnConnect = true
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

        Context("when UAA cannot be reached", func() {
            It("returns a 502 status code", func() {
                uaaClient.ErrorForUserByID = &url.Error{}
                handler.ServeHTTP(writer, request)

                Expect(writer.Code).To(Equal(http.StatusBadGateway))
            })
        })

        Context("when UAA cannot find the user", func() {
            It("returns a 410 status code", func() {
                uaaClient.ErrorForUserByID = uaa.Failure{}
                handler.ServeHTTP(writer, request)

                Expect(writer.Code).To(Equal(http.StatusGone))
            })
        })

        Context("when UAA causes some unknown error", func() {
            It("returns a 500 status code", func() {
                uaaClient.ErrorForUserByID = errors.New("Boom!")
                handler.ServeHTTP(writer, request)

                Expect(writer.Code).To(Equal(http.StatusInternalServerError))
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

    Context("when the request uses an expired auth token", func() {
        BeforeEach(func() {
            tokenHeader := map[string]interface{}{
                "alg": "RS256",
            }
            tokenClaims := map[string]interface{}{
                "jti":       "c5f6a266-5cf0-4ae2-9647-2615e7d28fa1",
                "client_id": "mister-client",
                "cid":       "mister-client",
                "exp":       1404281214,
            }
            token = BuildToken(tokenHeader, tokenClaims)

            requestBody, err := json.Marshal(map[string]string{
                "kind": "forgot_password",
                "text": "Try to remember your password next time",
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

        It("returns a 401 status code and error message", func() {
            handler.ServeHTTP(writer, request)

            Expect(writer.Code).To(Equal(http.StatusUnauthorized))

            parsed := map[string][]string{}
            err := json.Unmarshal(writer.Body.Bytes(), &parsed)
            if err != nil {
                panic(err)
            }

            Expect(parsed["errors"]).To(ContainElement("Authorization header is invalid: expired"))
        })
    })

    Context("when the request does not contain an auth token", func() {
        BeforeEach(func() {
            requestBody, err := json.Marshal(map[string]string{
                "kind": "forgot_password",
                "text": "Try to remember your password next time",
            })
            if err != nil {
                panic(err)
            }

            request, err = http.NewRequest("POST", "/users/user-123", bytes.NewReader(requestBody))
            if err != nil {
                panic(err)
            }
        })

        It("returns a 401 status code and error message", func() {
            handler.ServeHTTP(writer, request)

            Expect(writer.Code).To(Equal(http.StatusUnauthorized))

            parsed := map[string][]string{}
            err := json.Unmarshal(writer.Body.Bytes(), &parsed)
            if err != nil {
                panic(err)
            }

            Expect(parsed["errors"]).To(ContainElement("Authorization header is invalid: missing"))
        })
    })

    Context("when the auth token does not contain the correct scope", func() {
        BeforeEach(func() {
            tokenHeader := map[string]interface{}{
                "alg": "RS256",
            }
            tokenClaims := map[string]interface{}{
                "jti":       "c5f6a266-5cf0-4ae2-9647-2615e7d28fa1",
                "client_id": "mister-client",
                "cid":       "mister-client",
                "exp":       3404281214,
                "scope":     []string{"cloud_controller.admin"},
            }
            token = BuildToken(tokenHeader, tokenClaims)

            requestBody, err := json.Marshal(map[string]string{
                "kind": "forgot_password",
                "text": "Try to remember your password next time",
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

        It("returns a 403 status code and error message", func() {
            handler.ServeHTTP(writer, request)

            Expect(writer.Code).To(Equal(http.StatusForbidden))

            parsed := map[string][]string{}
            err := json.Unmarshal(writer.Body.Bytes(), &parsed)
            if err != nil {
                panic(err)
            }

            Expect(parsed["errors"]).To(ContainElement("You are not authorized to perform the requested action"))
        })
    })
})
