package handlers_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"

    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/cloudfoundry-incubator/notifications/web/handlers"

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
        var finder *FakeFinder
        var transaction *FakeDBConn
        var fakeRegistrar *FakeRegistrar

        BeforeEach(func() {
            errorWriter = &FakeErrorWriter{}

            finder = NewFakeFinder()
            finder.Clients["mister-client"] = models.Client{
                ID:          "mister-client",
                Description: "Health Monitor",
            }
            finder.Kinds["test_email|mister-client"] = models.Kind{
                ID:          "test_email",
                Description: "Instance Down",
                ClientID:    "mister-client",
            }

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
                "exp":       int64(3404281214),
                "scope":     []string{"notifications.write"},
            }
            token = BuildToken(tokenHeader, tokenClaims)

            request, err = http.NewRequest("POST", "/spaces/space-001", bytes.NewBuffer(body))
            if err != nil {
                panic(err)
            }
            request.Header.Set("Authorization", "Bearer "+token)

            transaction = &FakeDBConn{}

            fakeCourier = NewFakeCourier()
            fakeRegistrar = NewFakeRegistrar()
            handler = handlers.NewNotifySpace(handlers.NewNotify(fakeCourier, finder, fakeRegistrar), errorWriter)
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

                handler.Execute(writer, request, transaction)

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
                    "mister-client",
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
    })
})
