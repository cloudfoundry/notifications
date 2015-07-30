package campaigntypes_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/helpers"
	"github.com/cloudfoundry-incubator/notifications/web/v2/campaigntypes"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ShowHandler", func() {
	var (
		handler                 campaigntypes.ShowHandler
		campaignTypesCollection *fakes.CampaignTypesCollection
		context                 stack.Context
		writer                  *httptest.ResponseRecorder
		request                 *http.Request
		database                *fakes.Database
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

	It("returns information on a given campaign type", func() {
		campaignTypesCollection.GetCall.ReturnCampaignType = collections.CampaignType{
			ID:          "campaign-type-id-one",
			Name:        helpers.AddressOfString("first-campaign-type"),
			Description: helpers.AddressOfString("first-campaign-type-description"),
			Critical:    helpers.AddressOfBool(false),
			TemplateID:  helpers.AddressOfString(""),
			SenderID:    "some-sender-id",
		}
		var err error
		request, err = http.NewRequest("GET", "/senders/some-sender-id/campaign_types/campaign-type-id-one", nil)
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
		It("returns a 400 when the campaign type ID is an empty string", func() {
			campaignTypesCollection.GetCall.Err = collections.NewValidationError("missing campaign type id")

			var err error
			request, err = http.NewRequest("GET", "/senders/some-sender-id/campaign_types/", nil)
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusBadRequest))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"error": "missing campaign type id"
			}`))
		})

		It("returns a 400 when the sender ID is an empty string", func() {
			campaignTypesCollection.GetCall.Err = collections.NewValidationError("missing sender id")

			var err error
			request, err = http.NewRequest("GET", "/senders//campaign_types/some-campaign-type-id", nil)
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusBadRequest))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"error": "missing sender id"
			}`))
		})

		It("returns a 404 when the campaign type does not exist", func() {
			campaignTypesCollection.GetCall.Err = collections.NewNotFoundError("campaign type not found")

			var err error
			request, err = http.NewRequest("GET", "/senders/some-sender-id/campaign_types/missing-campaign-type-id", nil)
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusNotFound))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"error": "campaign type not found"
			}`))
		})

		It("returns a 404 when the sender does not exist", func() {
			campaignTypesCollection.GetCall.Err = collections.NewNotFoundError("sender not found")

			var err error
			request, err = http.NewRequest("GET", "/senders/missing-sender-id/campaign_types/some-campaign-type-id", nil)
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
			request, err = http.NewRequest("GET", "/senders/some-sender-id/campaign_types/some-campaign-type-id", nil)
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusInternalServerError))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"error": "BOOM!"
			}`))
		})
	})
})
