package handlers_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/web/handlers"
    "github.com/cloudfoundry-incubator/notifications/web/params"
    "github.com/cloudfoundry-incubator/notifications/web/services"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("UpdatePreferences", func() {
    Describe("Execute", func() {
        var handler handlers.UpdatePreferences
        var writer *httptest.ResponseRecorder
        var request *http.Request
        var updater *FakePreferenceUpdater
        var errorWriter *FakeErrorWriter
        var fakeDBConn *FakeDBConn

        BeforeEach(func() {
            fakeDBConn = &FakeDBConn{}
            builder := services.NewPreferencesBuilder()

            builder.Add("raptors", "door-opening", false)
            builder.Add("raptors", "feeding-time", true)
            builder.Add("dogs", "barking", false)

            body, err := json.MarshalIndent(builder, "", "  ")
            if err != nil {
                panic(err)
            }

            request, err = http.NewRequest("PATCH", "/user_preferences", bytes.NewBuffer(body))
            if err != nil {
                panic(err)
            }

            tokenHeader := map[string]interface{}{
                "alg": "FAST",
            }
            tokenClaims := map[string]interface{}{
                "user_id": "correct-user",
                "exp":     int64(3404281214),
            }

            token := BuildToken(tokenHeader, tokenClaims)

            request.Header.Set("Authorization", "Bearer "+token)

            errorWriter = NewFakeErrorWriter()
            updater = NewFakePreferenceUpdater()
            handler = handlers.NewUpdatePreferences(updater, errorWriter)
            writer = httptest.NewRecorder()
        })

        It("Passes The Correct Arguments to PreferenceUpdater Execute", func() {
            handler.Execute(writer, request, fakeDBConn)
            Expect(len(updater.ExecuteArguments)).To(Equal(2))

            preferencesArguments := updater.ExecuteArguments[0]

            Expect(preferencesArguments).To(ContainElement(models.Preference{
                ClientID: "raptors",
                KindID:   "door-opening",
                Email:    false,
            }))
            Expect(preferencesArguments).To(ContainElement(models.Preference{
                ClientID: "raptors",
                KindID:   "feeding-time",
                Email:    true,
            }))
            Expect(preferencesArguments).To(ContainElement(models.Preference{
                ClientID: "dogs",
                KindID:   "barking",
                Email:    false,
            }))

            Expect(updater.ExecuteArguments[1]).To(Equal("correct-user"))
        })

        It("Returns a 200 status code when the Preference object does not error", func() {
            handler.Execute(writer, request, fakeDBConn)

            Expect(writer.Code).To(Equal(http.StatusOK))
        })

        Context("when the JSON body cannot be parsed", func() {
            It("sends a params.ParseError to the error writer", func() {
                var err error
                request, err = http.NewRequest("PATCH", "/user_preferences", strings.NewReader(""))
                if err != nil {
                    panic(err)
                }

                handler.Execute(writer, request, fakeDBConn)

                Expect(errorWriter.Error).To(Equal(params.ParseError{}))
            })
        })
    })
})
