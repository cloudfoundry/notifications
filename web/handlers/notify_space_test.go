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
    "github.com/cloudfoundry-incubator/notifications/web/handlers/params"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("NotifySpace", func() {
    Describe("ServeHTTP", func() {
        var handler handlers.NotifySpace
        var writer *httptest.ResponseRecorder
        var request *http.Request
        var token string
        var fakeCourier *FakeCourier
        var errorWriter *FakeErrorWriter
        var clientsRepo *FakeClientsRepo
        var kindsRepo *FakeKindsRepo

        BeforeEach(func() {
            var err error

            errorWriter = &FakeErrorWriter{}

            conn := FakeDBConn{}
            clientsRepo = NewFakeClientsRepo()
            clientsRepo.Create(conn, models.Client{
                ID:          "mister-client",
                Description: "Health Monitor",
            })

            kindsRepo = NewFakeKindsRepo()
            kindsRepo.Create(conn, models.Kind{
                ID:          "test_email",
                Description: "Instance Down",
                ClientID:    "mister-client",
            })

            writer = httptest.NewRecorder()
            body, err := json.Marshal(map[string]string{
                "kind_id":  "test_email",
                "text":     "This is the plain text body of the email",
                "html":     "<p>This is the HTML Body of the email</p>",
                "subject":  "Your instance is down",
                "reply_to": "me@example.com",
            })
            if err != nil {
                panic(err)
            }

            tokenHeader := map[string]interface{}{
                "alg": "FAST",
            }
            tokenClaims := map[string]interface{}{
                "client_id": "mister-client",
                "exp":       3404281214,
                "scope":     []string{"notifications.write"},
            }
            token = BuildToken(tokenHeader, tokenClaims)

            request, err = http.NewRequest("POST", "/spaces/space-001", bytes.NewBuffer(body))
            if err != nil {
                panic(err)
            }
            request.Header.Set("Authorization", "Bearer "+token)

            fakeCourier = NewFakeCourier()
            handler = handlers.NewNotifySpace(fakeCourier, errorWriter, clientsRepo, kindsRepo)
        })

        Context("when the courier returns a successful response", func() {
            It("returns the JSON representation of the response", func() {
                fakeCourier.Responses = []postal.Response{
                    {
                        Status:         "delivered",
                        Recipient:      "user-123",
                        NotificationID: "abc-123",
                    },
                    {
                        Status:         "failed",
                        Recipient:      "user-456",
                        NotificationID: "abc-456",
                    },
                    {
                        Status:         "notfound",
                        Recipient:      "user-789",
                        NotificationID: "",
                    },
                    {
                        Status:         "noaddress",
                        Recipient:      "user-000",
                        NotificationID: "",
                    },
                }

                handler.ServeHTTP(writer, request)

                Expect(writer.Code).To(Equal(http.StatusOK))

                parsed := []map[string]string{}
                err := json.Unmarshal(writer.Body.Bytes(), &parsed)
                if err != nil {
                    panic(err)
                }

                Expect(len(parsed)).To(Equal(4))
                Expect(parsed).To(ContainElement(map[string]string{
                    "status":          "delivered",
                    "recipient":       "user-123",
                    "notification_id": "abc-123",
                }))
                Expect(parsed).To(ContainElement(map[string]string{
                    "status":          "failed",
                    "recipient":       "user-456",
                    "notification_id": "abc-456",
                }))
                Expect(parsed).To(ContainElement(map[string]string{
                    "status":          "notfound",
                    "recipient":       "user-789",
                    "notification_id": "",
                }))
                Expect(parsed).To(ContainElement(map[string]string{
                    "status":          "noaddress",
                    "recipient":       "user-000",
                    "notification_id": "",
                }))

                Expect(fakeCourier.DispatchArguments).To(Equal([]interface{}{
                    token,
                    postal.SpaceGUID("space-001"),
                    postal.Options{
                        ReplyTo:           "me@example.com",
                        Subject:           "Your instance is down",
                        KindDescription:   "Instance Down",
                        SourceDescription: "Health Monitor",
                        Text:              "This is the plain text body of the email",
                        HTML:              "<p>This is the HTML Body of the email</p>",
                        KindID:            "test_email",
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
                    request, err = http.NewRequest("POST", "/spaces/space-001", bytes.NewBuffer(body))
                    if err != nil {
                        panic(err)
                    }
                    request.Header.Set("Authorization", "Bearer "+token)

                    handler.ServeHTTP(writer, request)

                    Expect(errorWriter.Error).ToNot(BeNil())
                    validationErr := errorWriter.Error.(params.ValidationError)
                    Expect(validationErr.Errors()).To(ContainElement(`"kind_id" is a required field`))
                    Expect(validationErr.Errors()).To(ContainElement(`"text" or "html" fields must be supplied`))
                })

                It("returns a error response when params cannot be parsed", func() {
                    request, err := http.NewRequest("POST", "/spaces/space-001", strings.NewReader("this is not JSON"))
                    if err != nil {
                        panic(err)
                    }
                    request.Header.Set("Authorization", "Bearer "+token)

                    handler.ServeHTTP(writer, request)

                    Expect(errorWriter.Error).To(Equal(params.ParseError{}))
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
                            postal.SpaceGUID("space-001"),
                            postal.Options{
                                ReplyTo:           "me@example.com",
                                Subject:           "Your instance is down",
                                KindDescription:   "Instance Down",
                                SourceDescription: "",
                                Text:              "This is the plain text body of the email",
                                HTML:              "<p>This is the HTML Body of the email</p>",
                                KindID:            "test_email",
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
                            postal.SpaceGUID("space-001"),
                            postal.Options{
                                ReplyTo:           "me@example.com",
                                Subject:           "Your instance is down",
                                KindDescription:   "",
                                SourceDescription: "Health Monitor",
                                Text:              "This is the plain text body of the email",
                                HTML:              "<p>This is the HTML Body of the email</p>",
                                KindID:            "test_email",
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
})
