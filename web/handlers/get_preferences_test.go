package handlers_test

import (
    "bytes"
    "encoding/json"
    "errors"
    "net/http"
    "net/http/httptest"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/fakes"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/web/handlers"
    "github.com/cloudfoundry-incubator/notifications/web/services"
    "github.com/dgrijalva/jwt-go"
    "github.com/ryanmoran/stack"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("GetPreferences", func() {
    var handler handlers.GetPreferences
    var writer *httptest.ResponseRecorder
    var request *http.Request
    var preferencesFinder *fakes.FakePreferencesFinder
    var errorWriter *fakes.ErrorWriter
    var builder services.PreferencesBuilder
    var context stack.Context
    var TRUE = true
    var FALSE = false

    BeforeEach(func() {
        errorWriter = fakes.NewErrorWriter()

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

        token, err := jwt.Parse(fakes.BuildToken(tokenHeader, tokenClaims), func(token *jwt.Token) ([]byte, error) {
            return []byte(config.UAAPublicKey), nil
        })
        context = stack.NewContext()
        context.Set("token", token)

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
        builder.GlobalUnsubscribe = true

        preferencesFinder = fakes.NewFakePreferencesFinder(builder)
        handler = handlers.NewGetPreferences(preferencesFinder, errorWriter)
    })

    It("Passes the proper user guid into execute", func() {
        handler.ServeHTTP(writer, request, context)

        Expect(preferencesFinder.UserGUID).To(Equal("correct-user"))
    })

    It("Returns a proper JSON response when the Preference object does not error", func() {
        handler.ServeHTTP(writer, request, context)

        Expect(writer.Code).To(Equal(http.StatusOK))

        parsed := services.PreferencesBuilder{}
        err := json.Unmarshal(writer.Body.Bytes(), &parsed)
        if err != nil {
            panic(err)
        }

        Expect(parsed.GlobalUnsubscribe).To(BeTrue())
        Expect(parsed.Clients["raptorClient"]["hungry-kind"].Email).To(Equal(&FALSE))
        Expect(parsed.Clients["raptorClient"]["hungry-kind"].Count).To(Equal(0))
        Expect(parsed.Clients["starWarsClient"]["vader-kind"].Email).To(Equal(&TRUE))
        Expect(parsed.Clients["starWarsClient"]["vader-kind"].Count).To(Equal(0))
    })

    Context("when there is an error returned from the finder", func() {
        It("writes the error to the error writer", func() {
            preferencesFinder.FindError = errors.New("boom!")
            handler.ServeHTTP(writer, request, context)
            Expect(errorWriter.Error).To(Equal(preferencesFinder.FindError))
        })
    })
})
