package handlers_test

import (
    "bytes"
    "encoding/json"
    "errors"
    "net/http"
    "net/http/httptest"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/cloudfoundry-incubator/notifications/web/handlers"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("NotifyUser", func() {
    var handler handlers.NotifyUser
    var writer *httptest.ResponseRecorder
    var request *http.Request
    var token string
    var fakeCourier *FakeCourier
    var errorWriter *FakeErrorWriter
    var clientsRepo *FakeClientsRepo
    var kindsRepo *FakeKindsRepo

    BeforeEach(func() {
        tokenHeader := map[string]interface{}{
            "alg": "FAST",
        }
        tokenClaims := map[string]interface{}{
            "client_id": "mister-client",
            "exp":       3404281214,
            "scope":     []string{"notifications.write"},
        }
        token = BuildToken(tokenHeader, tokenClaims)
        writer = httptest.NewRecorder()

        fakeCourier = NewFakeCourier()
        errorWriter = &FakeErrorWriter{}

        conn := FakeDBConn{}
        clientsRepo = NewFakeClientsRepo()
        clientsRepo.Create(conn, models.Client{
            ID:          "mister-client",
            Description: "Login System",
        })

        kindsRepo = NewFakeKindsRepo()
        kindsRepo.Create(conn, models.Kind{
            ID:          "forgot_password",
            Description: "Password Reminder",
            ClientID:    "mister-client",
        })

        handler = handlers.NewNotifyUser(fakeCourier, errorWriter, clientsRepo, kindsRepo)

        requestBody, err := json.Marshal(map[string]string{
            "kind_id":  "forgot_password",
            "subject":  "Forgot password request",
            "text":     "Please reset your password by clicking on this link...",
            "html":     "<p>Please reset your password by clicking on this link...</p>",
            "reply_to": "me@example.com",
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

    Context("when the courier returns a successful response", func() {
        It("returns the JSON representation of the response", func() {
            fakeCourier.Responses = []postal.Response{
                {
                    Status:         "delivered",
                    Recipient:      "user-123",
                    NotificationID: "abc-123",
                },
            }

            handler.ServeHTTP(writer, request)

            Expect(writer.Code).To(Equal(http.StatusOK))

            parsed := []map[string]string{}
            err := json.Unmarshal(writer.Body.Bytes(), &parsed)
            if err != nil {
                panic(err)
            }

            Expect(len(parsed)).To(Equal(1))
            Expect(parsed).To(ContainElement(map[string]string{
                "status":          "delivered",
                "recipient":       "user-123",
                "notification_id": "abc-123",
            }))
            Expect(fakeCourier.DispatchArguments).To(Equal([]interface{}{
                token,
                postal.UserGUID("user-123"),
                postal.Options{
                    ReplyTo:           "me@example.com",
                    Subject:           "Forgot password request",
                    KindDescription:   "Password Reminder",
                    SourceDescription: "Login System",
                    Text:              "Please reset your password by clicking on this link...",
                    HTML:              "<p>Please reset your password by clicking on this link...</p>",
                    KindID:            "forgot_password",
                },
            }))
        })
    })

    Context("failure cases", func() {
        Context("when validating params", func() {
            It("returns a error response when params are missing", func() {
                body, err := json.Marshal(map[string]string{
                    "subject": "Your instance is down",
                })
                if err != nil {
                    panic(err)
                }
                request, err = http.NewRequest("POST", "/users/user-001", bytes.NewBuffer(body))
                if err != nil {
                    panic(err)
                }
                request.Header.Set("Authorization", "Bearer "+token)

                handler.ServeHTTP(writer, request)

                Expect(errorWriter.Error).ToNot(BeNil())
                validationErr := errorWriter.Error.(handlers.ParamsValidationError)
                Expect(validationErr.Errors()).To(ContainElement(`"kind_id" is a required field`))
                Expect(validationErr.Errors()).To(ContainElement(`"text" or "html" fields must be supplied`))
            })

            It("returns a error response when params cannot be parsed", func() {
                request, err := http.NewRequest("POST", "/users/user-001", strings.NewReader("this is not JSON"))
                if err != nil {
                    panic(err)
                }
                request.Header.Set("Authorization", "Bearer "+token)

                handler.ServeHTTP(writer, request)

                Expect(errorWriter.Error).To(Equal(handlers.ParamsParseError{}))
            })
        })

        Context("when the courier returns errors", func() {
            It("delegates to the errorWriter", func() {
                fakeCourier.Error = errors.New("BOOM!")

                handler.ServeHTTP(writer, request)

                Expect(errorWriter.Error).To(Equal(errors.New("BOOM!")))
            })
        })

        Context("when the client repo return errors", func() {
            Context("when the error is a RecordNotFound", func() {
                It("does not populate the SourceDescription", func() {
                    clientsRepo.Clients = map[string]models.Client{}

                    handler.ServeHTTP(writer, request)

                    Expect(errorWriter.Error).To(BeNil())
                    Expect(fakeCourier.DispatchArguments).To(Equal([]interface{}{
                        token,
                        postal.UserGUID("user-123"),
                        postal.Options{
                            ReplyTo:           "me@example.com",
                            Subject:           "Forgot password request",
                            KindDescription:   "Password Reminder",
                            SourceDescription: "",
                            Text:              "Please reset your password by clicking on this link...",
                            HTML:              "<p>Please reset your password by clicking on this link...</p>",
                            KindID:            "forgot_password",
                        },
                    }))
                })
            })

            Context("when the error is not a RecordNotFound", func() {
                It("delegates to the errorWriter", func() {
                    clientsRepo.FindError = errors.New("BOOM!")

                    handler.ServeHTTP(writer, request)

                    Expect(errorWriter.Error).To(Equal(errors.New("BOOM!")))
                })
            })
        })

        Context("when the kind repo return errors", func() {
            Context("when the error is a RecordNotFound", func() {
                It("does not populate the KindDescription", func() {
                    kindsRepo.Kinds = map[string]models.Kind{}

                    handler.ServeHTTP(writer, request)

                    Expect(errorWriter.Error).To(BeNil())
                    Expect(fakeCourier.DispatchArguments).To(Equal([]interface{}{
                        token,
                        postal.UserGUID("user-123"),
                        postal.Options{
                            ReplyTo:           "me@example.com",
                            Subject:           "Forgot password request",
                            KindDescription:   "",
                            SourceDescription: "Login System",
                            Text:              "Please reset your password by clicking on this link...",
                            HTML:              "<p>Please reset your password by clicking on this link...</p>",
                            KindID:            "forgot_password",
                        },
                    }))
                })
            })

            Context("when the error is not a RecordNotFound", func() {
                It("delegates to the errorWriter", func() {
                    kindsRepo.FindError = errors.New("BOOM!")

                    handler.ServeHTTP(writer, request)

                    Expect(errorWriter.Error).To(Equal(errors.New("BOOM!")))
                })
            })
        })
    })
})
