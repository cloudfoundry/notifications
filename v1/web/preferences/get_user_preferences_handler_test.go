package preferences_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/testing/fakes"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/cloudfoundry-incubator/notifications/v1/web/preferences"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetUserPreferencesHandler", func() {
	var (
		handler           preferences.GetUserPreferencesHandler
		writer            *httptest.ResponseRecorder
		request           *http.Request
		preferencesFinder *fakes.PreferencesFinder
		errorWriter       *fakes.ErrorWriter
		builder           services.PreferencesBuilder
		context           stack.Context
		database          *fakes.Database
	)

	BeforeEach(func() {
		errorWriter = fakes.NewErrorWriter()

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

		database = fakes.NewDatabase()
		context = stack.NewContext()
		context.Set("database", database)

		preferencesFinder = fakes.NewPreferencesFinder(builder)
		handler = preferences.NewGetUserPreferencesHandler(preferencesFinder, errorWriter)
	})

	Context("when a client is making a request for an arbitrary user", func() {
		It("passes the proper user guid to the finder", func() {
			handler.ServeHTTP(writer, request, context)
			Expect(preferencesFinder.FindCall.Arguments).To(ConsistOf([]interface{}{database, "af02af02-af02-af02-af02-af02af02af02"}))
		})

		It("returns a proper JSON response when the Preference object does not error", func() {
			handler.ServeHTTP(writer, request, context)

			Expect(writer.Code).To(Equal(http.StatusOK))

			Expect(string(writer.Body.Bytes())).To(Equal(`{"global_unsubscribe":false,"clients":{"raptorClient":{"hungry-kind":{"count":0,"email":false,"kind_description":"hungry-kind","source_description":"raptorClient"}},"starWarsClient":{"vader-kind":{"count":0,"email":true,"kind_description":"vader-kind","source_description":"starWarsClient"}}}}`))
		})

		Context("when the finder returns an error", func() {
			It("writes the error to the error writer", func() {
				preferencesFinder.FindCall.Error = errors.New("wow!!")
				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.Error).To(Equal(preferencesFinder.FindCall.Error))
			})
		})
	})
})
