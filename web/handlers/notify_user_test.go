package handlers_test

import (
    "bytes"
    "encoding/json"
    "errors"
    "net/http"
    "net/http/httptest"

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
        handler = handlers.NewNotifyUser(fakeCourier)

        requestBody, err := json.Marshal(map[string]string{
            "kind":               "forgot_password",
            "kind_description":   "Password reminder",
            "source_description": "Login system",
            "text":               "Please reset your password by clicking on this link...",
            "html":               "<p>Please reset your password by clicking on this link...</p>",
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
        })
    })

    Context("when validating params", func() {
        It("returns a error response when params are missing", func() {
            body, err := json.Marshal(map[string]string{
                "subject":            "Your instance is down",
                "source_description": "MySQL Service",
                "kind_description":   "Instance Alert",
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

            parsed := map[string][]string{}
            err = json.Unmarshal(writer.Body.Bytes(), &parsed)
            if err != nil {
                panic(err)
            }

            Expect(parsed["errors"]).To(ContainElement(`"kind" is a required field`))
            Expect(parsed["errors"]).To(ContainElement(`"text" or "html" fields must be supplied`))
        })
    })

    Context("when the courier returns errors", func() {
        It("returns a 502 when CloudController fails to respond", func() {
            fakeCourier.Error = postal.CCDownError("BOOM!")

            handler.ServeHTTP(writer, request)

            Expect(writer.Code).To(Equal(http.StatusBadGateway))

            body := make(map[string]interface{})
            err := json.Unmarshal(writer.Body.Bytes(), &body)
            if err != nil {
                panic(err)
            }

            Expect(body["errors"]).To(ContainElement("Cloud Controller is unavailable"))
        })

        It("returns a 502 when UAA fails to respond", func() {
            fakeCourier.Error = postal.UAADownError("BOOM!")

            handler.ServeHTTP(writer, request)

            Expect(writer.Code).To(Equal(http.StatusBadGateway))

            body := make(map[string]interface{})
            err := json.Unmarshal(writer.Body.Bytes(), &body)
            if err != nil {
                panic(err)
            }

            Expect(body["errors"]).To(ContainElement("UAA is unavailable"))
        })

        It("returns a 502 when UAA fails for unknown reasons", func() {
            fakeCourier.Error = postal.UAAGenericError("UAA Unknown Error: BOOM!")

            handler.ServeHTTP(writer, request)

            Expect(writer.Code).To(Equal(http.StatusBadGateway))

            body := make(map[string]interface{})
            err := json.Unmarshal(writer.Body.Bytes(), &body)
            if err != nil {
                panic(err)
            }

            Expect(body["errors"]).To(ContainElement("UAA Unknown Error: BOOM!"))
        })

        It("returns a 500 when the is a template loading error", func() {
            fakeCourier.Error = postal.TemplateLoadError("BOOM!")

            handler.ServeHTTP(writer, request)

            Expect(writer.Code).To(Equal(http.StatusInternalServerError))

            body := make(map[string]interface{})
            err := json.Unmarshal(writer.Body.Bytes(), &body)
            if err != nil {
                panic(err)
            }

            Expect(body["errors"]).To(ContainElement("An email template could not be loaded"))
        })

        It("panics for unknown errors", func() {
            fakeCourier.Error = errors.New("BOOM!")

            Expect(func() {
                handler.ServeHTTP(writer, request)
            }).To(Panic())
        })
    })
})
