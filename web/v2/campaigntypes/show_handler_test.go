package campaigntypes_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/web/v2/campaigntypes"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ShowHandler", func() {
	var (
		handler                     campaigntypes.ShowHandler
		campaignTypesCollection *fakes.CampaignTypesCollection
		context                     stack.Context
		writer                      *httptest.ResponseRecorder
		request                     *http.Request
		database                    *fakes.Database
	)

	BeforeEach(func() {
		context = stack.NewContext()

		context.Set("client_id", "some-client-id")

		database = fakes.NewDatabase()
		context.Set("database", database)

		writer = httptest.NewRecorder()

		campaignTypesCollection = fakes.NewCampaignTypesCollection()

		handler = campaigntypes.NewShowHandler(campaignTypesCollection)
	})

	It("returns information on a given notification type", func() {
		campaignTypesCollection.GetCall.ReturnCampaignType = collections.CampaignType{
			ID:          "campaign-type-id-one",
			Name:        "first-campaign-type",
			Description: "first-campaign-type-description",
			Critical:    false,
			TemplateID:  "",
			SenderID:    "some-sender-id",
		}
		var err error
		request, err = http.NewRequest("GET", "/senders/some-sender-id/notification_types/campaign-type-id-one", nil)
		Expect(err).NotTo(HaveOccurred())

		handler.ServeHTTP(writer, request, context)

		Expect(writer.Code).To(Equal(http.StatusOK))
		Expect(writer.Body.String()).To(MatchJSON(`{
			"id": "campaign-type-id-one",
			"name": "first-campaign-type",
			"description": "first-campaign-type-description",
			"critical": false,
			"template_id": ""
		}`))
	})

	Context("failure cases", func() {
		It("returns a 400 when the notification type ID is an empty string", func() {
			campaignTypesCollection.GetCall.Err = collections.NewValidationError("missing notification type id")

			var err error
			request, err = http.NewRequest("GET", "/senders/some-sender-id/notification_types/", nil)
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusBadRequest))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"error": "missing notification type id"
			}`))
		})

		It("returns a 400 when the sender ID is an empty string", func() {
			campaignTypesCollection.GetCall.Err = collections.NewValidationError("missing sender id")

			var err error
			request, err = http.NewRequest("GET", "/senders//notification_types/some-campaign-type-id", nil)
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusBadRequest))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"error": "missing sender id"
			}`))
		})

		It("returns a 404 when the notification type does not exist", func() {
			campaignTypesCollection.GetCall.Err = collections.NewNotFoundError("notification type not found")

			var err error
			request, err = http.NewRequest("GET", "/senders/some-sender-id/notification_types/missing-campaign-type-id", nil)
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusNotFound))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"error": "notification type not found"
			}`))
		})

		It("returns a 404 when the sender does not exist", func() {
			campaignTypesCollection.GetCall.Err = collections.NewNotFoundError("sender not found")

			var err error
			request, err = http.NewRequest("GET", "/senders/missing-sender-id/notification_types/some-campaign-type-id", nil)
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusNotFound))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"error": "sender not found"
			}`))
		})

		It("returns a 500 when the collection indicates a system error", func() {
			campaignTypesCollection.GetCall.Err = errors.New("BOOM!")

			var err error
			request, err = http.NewRequest("GET", "/senders/some-sender-id/notification_types/some-campaign-type-id", nil)
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusInternalServerError))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"error": "BOOM!"
			}`))
		})
	})
})
