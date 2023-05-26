package preferences_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/cloudfoundry-incubator/notifications/v1/web/preferences"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetUserPreferencesHandler", func() {
	var (
		handler           preferences.GetUserPreferencesHandler
		writer            *httptest.ResponseRecorder
		request           *http.Request
		preferencesFinder *mocks.PreferencesFinder
		errorWriter       *mocks.ErrorWriter
		builder           services.PreferencesBuilder
		context           stack.Context
		database          *mocks.Database
	)

	BeforeEach(func() {
		errorWriter = mocks.NewErrorWriter()

		writer = httptest.NewRecorder()
		body, err := json.Marshal(map[string]string{
			"I think this request is empty": "maybe",
		})
		Expect(err).NotTo(HaveOccurred())

		request, err = http.NewRequest("GET", "/user_preferences/af02af02-af02-af02-af02-af02af02af02", bytes.NewBuffer(body))
		Expect(err).NotTo(HaveOccurred())

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

		database = mocks.NewDatabase()
		context = stack.NewContext()
		context.Set("database", database)

		preferencesFinder = mocks.NewPreferencesFinder()
		preferencesFinder.FindCall.Returns.PreferencesBuilder = builder

		handler = preferences.NewGetUserPreferencesHandler(preferencesFinder, errorWriter)
	})

	Context("when a client is making a request for an arbitrary user", func() {
		It("passes the proper user guid to the finder", func() {
			handler.ServeHTTP(writer, request, context)
			Expect(preferencesFinder.FindCall.Receives.Database).To(Equal(database))
			Expect(preferencesFinder.FindCall.Receives.UserGUID).To(Equal("af02af02-af02-af02-af02-af02af02af02"))
		})

		It("returns a proper JSON response when the Preference object does not error", func() {
			handler.ServeHTTP(writer, request, context)

			Expect(writer.Code).To(Equal(http.StatusOK))

			Expect(writer.Body).To(MatchJSON(`{
				"global_unsubscribe":false,
				"clients":{
					"raptorClient":{
						"hungry-kind":{
							"email":false,
							"kind_description":"hungry-kind",
							"source_description":"raptorClient"
						}
					},
					"starWarsClient":{
						"vader-kind":{
							"email":true,
							"kind_description":"vader-kind",
							"source_description":"starWarsClient"
						}
					}
				}
			}`))
		})

		Context("when the finder returns an error", func() {
			It("writes the error to the error writer", func() {
				preferencesFinder.FindCall.Returns.Error = errors.New("wow!!")
				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.WriteCall.Receives.Error).To(Equal(preferencesFinder.FindCall.Returns.Error))
			})
		})
	})
})
