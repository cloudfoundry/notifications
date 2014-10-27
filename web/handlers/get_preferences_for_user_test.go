package handlers_test

import (
    "bytes"
    "encoding/json"
    "errors"
    "net/http"
    "net/http/httptest"

    "github.com/cloudfoundry-incubator/notifications/fakes"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/web/handlers"
    "github.com/cloudfoundry-incubator/notifications/web/services"
    "github.com/ryanmoran/stack"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("GetPreferencesForUser", func() {
    var handler handlers.GetPreferencesForUser
    var writer *httptest.ResponseRecorder
    var request *http.Request
    var preferencesFinder *fakes.FakePreferencesFinder
    var errorWriter *fakes.FakeErrorWriter
    var builder services.PreferencesBuilder
    var context stack.Context

    BeforeEach(func() {
        errorWriter = &fakes.FakeErrorWriter{}

        writer = httptest.NewRecorder()
        body, err := json.Marshal(map[string]string{
            "I think this request is empty": "maybe",
        })
        if err != nil {
            panic(err)
        }

        request, err = http.NewRequest("GET", "/user_preferences/af02af02-af02-af02-af02-af02af02af02", bytes.NewBuffer(body))
        if err != nil {
            panic(err)
        }

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

        preferencesFinder = fakes.NewFakePreferencesFinder(builder)
        handler = handlers.NewGetPreferencesForUser(preferencesFinder, errorWriter)
    })

    Context("when a client is making a request for an arbitrary user", func() {

        It("Passes the proper user guid to the finder", func() {
            handler.ServeHTTP(writer, request, context)
            Expect(preferencesFinder.UserGUID).To(Equal("af02af02-af02-af02-af02-af02af02af02"))
        })

        It("Returns a proper JSON response when the Preference object does not error", func() {
            handler.ServeHTTP(writer, request, context)

            Expect(writer.Code).To(Equal(http.StatusOK))

            Expect(string(writer.Body.Bytes())).To(Equal(`{"global_unsubscribe":false,"clients":{"raptorClient":{"hungry-kind":{"count":0,"email":false,"kind_description":"hungry-kind","source_description":"raptorClient"}},"starWarsClient":{"vader-kind":{"count":0,"email":true,"kind_description":"vader-kind","source_description":"starWarsClient"}}}}`))
        })

        Context("when the finder returns an error", func() {
            It("writes the error to the error writer", func() {
                preferencesFinder.FindError = errors.New("wow!!")
                handler.ServeHTTP(writer, request, context)
                Expect(errorWriter.Error).To(Equal(preferencesFinder.FindError))
            })
        })
    })
})
