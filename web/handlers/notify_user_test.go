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

var _ = Describe("NotifyUser", func() {
    Context("ServeHTTP", func() {
        var handler handlers.NotifyUser
        var writer *httptest.ResponseRecorder
        var request *http.Request
        var token string
        var fakeCourier *FakeCourier
        var errorWriter *FakeErrorWriter
        var finder *FakeFinder
        var fakeRegistrar *FakeRegistrar
        var transaction *FakeDBConn

        BeforeEach(func() {
            errorWriter = &FakeErrorWriter{}

            finder = NewFakeFinder()
            finder.Clients["mister-client"] = models.Client{
                ID:          "mister-client",
                Description: "Login System",
            }
            finder.Kinds["forgot_password|mister-client"] = models.Kind{
                ID:          "forgot_password",
                Description: "Password Reminder",
                ClientID:    "mister-client",
            }

            writer = httptest.NewRecorder()
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

            tokenHeader := map[string]interface{}{
                "alg": "FAST",
            }
            tokenClaims := map[string]interface{}{
                "client_id": "mister-client",
                "exp":       int64(3404281214),
                "scope":     []string{"notifications.write"},
            }
            token = BuildToken(tokenHeader, tokenClaims)

            request, err = http.NewRequest("POST", "/users/user-123", bytes.NewReader(requestBody))
            if err != nil {
                panic(err)
            }
            request.Header.Set("Authorization", "Bearer "+token)

            transaction = &FakeDBConn{}
            fakeCourier = NewFakeCourier()
            fakeRegistrar = NewFakeRegistrar()
            handler = handlers.NewNotifyUser(handlers.NewNotify(fakeCourier, finder, fakeRegistrar), errorWriter)
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

                handler.Execute(writer, request, transaction)

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
                    "mister-client",
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
    })
})
