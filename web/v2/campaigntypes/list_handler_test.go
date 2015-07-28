package campaigntypes_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/web/v2/campaigntypes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ryanmoran/stack"
)

var _ = Describe("ListHandler", func() {
	var (
		handler                     campaigntypes.ListHandler
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

		handler = campaigntypes.NewListHandler(campaignTypesCollection)
	})

	It("returns a list of notification types", func() {
		campaignTypesCollection.ListCall.ReturnCampaignTypeList = []collections.CampaignType{
			{
				ID:          "campaign-type-id-one",
				Name:        "first-campaign-type",
				Description: "first-campaign-type-description",
				Critical:    false,
				TemplateID:  "",
				SenderID:    "some-sender-id",
			},
			{
				ID:          "campaign-type-id-two",
				Name:        "second-campaign-type",
				Description: "second-campaign-type-description",
				Critical:    true,
				TemplateID:  "",
				SenderID:    "some-sender-id",
			},
		}

		var err error
		request, err = http.NewRequest("GET", "/senders/some-sender-id/campaign_types", nil)
		Expect(err).NotTo(HaveOccurred())

		handler.ServeHTTP(writer, request, context)

		Expect(writer.Code).To(Equal(http.StatusOK))
		Expect(writer.Body.String()).To(MatchJSON(`{
			"campaign_types": [
				{
					"id": "campaign-type-id-one",
					"name": "first-campaign-type",
					"description": "first-campaign-type-description",
					"critical": false,
					"template_id": ""
				},
				{
					"id": "campaign-type-id-two",
					"name": "second-campaign-type",
					"description": "second-campaign-type-description",
					"critical": true,
					"template_id": ""
				}
			]
		}`))
	})

	It("returns an empty list of notification types if the table has no records", func() {
		var err error
		request, err = http.NewRequest("GET", "/senders/some-sender-id/campaign_types", nil)
		Expect(err).NotTo(HaveOccurred())

		handler.ServeHTTP(writer, request, context)

		Expect(writer.Code).To(Equal(http.StatusOK))
		Expect(writer.Body.String()).To(MatchJSON(`{
			"campaign_types": []
		}`))
	})

	Context("failure cases", func() {
		It("returns a 404 when the sender does not exist", func() {
			campaignTypesCollection.ListCall.Err = collections.NotFoundError{
				Err: errors.New("sender not found"),
			}

			var err error
			request, err = http.NewRequest("GET", "/senders/non-existent-sender-id/campaign_types", nil)
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusNotFound))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"error": "sender not found"
			}`))
		})

		It("returns a 500 when the collection indicates a system error", func() {
			campaignTypesCollection.ListCall.Err = errors.New("BOOM!")

			var err error
			request, err = http.NewRequest("GET", "/senders/some-sender-id/campaign_types", nil)
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusInternalServerError))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"error": "BOOM!"
			}`))
		})
	})
})
