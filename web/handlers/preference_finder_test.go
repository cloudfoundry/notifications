package handlers_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"

    "github.com/cloudfoundry-incubator/notifications/web/handlers"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("PreferenceFinder", func() {

    var preferenceFinder handlers.PreferenceFinder
    var writer *httptest.ResponseRecorder
    var request *http.Request
    var preference *FakePreference
    var errorWriter *FakeErrorWriter
    var preferencesMap handlers.NotificationPreferences

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
            "exp":     3404281214,
            "scope":   []string{"notifications.write"},
        }

        token := BuildToken(tokenHeader, tokenClaims)

        request, err = http.NewRequest("POST", "/preferences", bytes.NewBuffer(body))
        if err != nil {
            panic(err)
        }

        request.Header.Set("Authorization", "Bearer "+token)

        preferencesMap = map[string]map[string]map[string]string{
            "raptorClient": map[string]map[string]string{
                "hungry-kind": map[string]string{"email": "false"},
            },

            "starWarsClient": map[string]map[string]string{
                "vader-kind": map[string]string{"email": "true"},
            },
        }

        preference = NewFakePreference(preferencesMap)

        preferenceFinder = handlers.NewPreferenceFinder(preference, errorWriter)
    })

    It("Passes the proper user guid into execute", func() {
        preferenceFinder.ServeHTTP(writer, request)

        Expect(preference.UserGUID).To(Equal("correct-user"))

    })

    It("Returns a proper JSON response when the Preference object does not error", func() {

        preferenceFinder.ServeHTTP(writer, request)

        Expect(writer.Code).To(Equal(http.StatusOK))

        parsed := handlers.NotificationPreferences{}
        err := json.Unmarshal(writer.Body.Bytes(), &parsed)
        if err != nil {
            panic(err)
        }

        Expect(parsed).To(Equal(preferencesMap))

    })

    Context("when there is a database error", func() {
        It("panics", func() {
            preference.ExecuteErrors = true
            Expect(func() { preferenceFinder.ServeHTTP(writer, request) }).To(Panic())
        })
    })
})
