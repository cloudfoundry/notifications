package handlers_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"

    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/web/handlers"
    "github.com/cloudfoundry-incubator/notifications/web/services"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("GetPreferences", func() {
    var handler handlers.GetPreferences
    var writer *httptest.ResponseRecorder
    var request *http.Request
    var preferencesFinder *FakePreferencesFinder
    var errorWriter *FakeErrorWriter
    var builder services.PreferencesBuilder

    BeforeEach(func() {
        errorWriter = &FakeErrorWriter{}

        writer = httptest.NewRecorder()
        body, err := json.Marshal(map[string]string{
            "I think this request is empty": "maybe",
        })
        if err != nil {
            panic(err)
        }

        tokenHeader := map[string]interface{}{
            "alg": "FAST",
        }
        tokenClaims := map[string]interface{}{
            "user_id": "correct-user",
            "exp":     int64(3404281214),
            "scope":   []string{"notification_preferences.read"},
        }

        request, err = http.NewRequest("GET", "/user_preferences", bytes.NewBuffer(body))
        if err != nil {
            panic(err)
        }

        request.Header.Set("Authorization", "Bearer "+BuildToken(tokenHeader, tokenClaims))

        builder = services.NewPreferencesBuilder()
        builder.Add(models.Preference{
            ClientID: "raptorClient",
            KindID:   "hungry-kind",
            Email:    false,
        })
        builder.Add(models.Preference{
            ClientID: "starWarsClient",
            KindID:   "vader-kind",
            Email:    true,
        })

        preferencesFinder = NewFakePreferencesFinder(builder)
        handler = handlers.NewGetPreferences(preferencesFinder, errorWriter)
    })

    It("Passes the proper user guid into execute", func() {
        handler.ServeHTTP(writer, request)

        Expect(preferencesFinder.UserGUID).To(Equal("correct-user"))
    })

    It("Returns a proper JSON response when the Preference object does not error", func() {
        handler.ServeHTTP(writer, request)

        Expect(writer.Code).To(Equal(http.StatusOK))

        parsed := services.PreferencesBuilder{}
        err := json.Unmarshal(writer.Body.Bytes(), &parsed)
        if err != nil {
            panic(err)
        }

        Expect(parsed).To(Equal(builder))
    })

    Context("when there is a database error", func() {
        It("panics", func() {
            preferencesFinder.FindErrors = true

            Expect(func() {
                handler.ServeHTTP(writer, request)
            }).To(Panic())
        })
    })
})
