package campaigntypes_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/web/campaigntypes"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ShowHandler", func() {
	var (
		handler                 campaigntypes.ShowHandler
		campaignTypesCollection *mocks.CampaignTypesCollection
		context                 stack.Context
		writer                  *httptest.ResponseRecorder
		request                 *http.Request
		database                *mocks.Database
		conn                    *mocks.Connection
	)

	BeforeEach(func() {
		context = stack.NewContext()

		context.Set("client_id", "some-client-id")

		conn = mocks.NewConnection()
		database = mocks.NewDatabase()
		database.ConnectionCall.Returns.Connection = conn
		context.Set("database", database)

		writer = httptest.NewRecorder()

		campaignTypesCollection = mocks.NewCampaignTypesCollection()

		handler = campaigntypes.NewShowHandler(campaignTypesCollection)
	})

	It("returns information on a given campaign type", func() {
		campaignTypesCollection.GetCall.Returns.CampaignType = collections.CampaignType{
			ID:          "campaign-type-id-one",
			Name:        "first-campaign-type",
			Description: "first-campaign-type-description",
			Critical:    true,
			TemplateID:  "template-id",
			SenderID:    "some-sender-id",
		}
		var err error
		request, err = http.NewRequest("GET", "/senders/some-sender-id/campaign_types/campaign-type-id-one", nil)
		Expect(err).NotTo(HaveOccurred())

		handler.ServeHTTP(writer, request, context)

		Expect(campaignTypesCollection.GetCall.Receives.Conn).To(Equal(conn))
		Expect(campaignTypesCollection.GetCall.Receives.CampaignTypeID).To(Equal("campaign-type-id-one"))
		Expect(campaignTypesCollection.GetCall.Receives.SenderID).To(Equal("some-sender-id"))
		Expect(campaignTypesCollection.GetCall.Receives.ClientID).To(Equal("some-client-id"))

		Expect(writer.Code).To(Equal(http.StatusOK))
		Expect(writer.Body.String()).To(MatchJSON(`{
			"id": "campaign-type-id-one",
			"name": "first-campaign-type",
			"description": "first-campaign-type-description",
			"critical": true,
			"template_id": "template-id",
			"_links": {
				"self": {
					"href": "/campaign_types/campaign-type-id-one"
				}
			}
		}`))
	})

	Context("failure cases", func() {
		It("returns a 301 when the campaign type id is missing", func() {
			var err error
			request, err = http.NewRequest("GET", "/senders/some-sender-id/campaign_types/", nil)
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusMovedPermanently))
			headers := writer.Header()
			Expect(headers.Get("Location")).To(Equal("/senders/some-sender-id/campaign_types"))
			Expect(writer.Body.String()).To(BeEmpty())
		})

		It("returns a 400 when the sender ID is an empty string", func() {
			var err error
			request, err = http.NewRequest("GET", "/senders//campaign_types/some-campaign-type-id", nil)
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusBadRequest))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"errors": ["missing sender id"]
			}`))
		})

		It("returns a 401 when the client id is missing", func() {
			var err error
			request, err = http.NewRequest("GET", "/senders/some-sender-id/campaign_types/missing-campaign-type-id", nil)
			Expect(err).NotTo(HaveOccurred())

			context.Set("client_id", "")

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusUnauthorized))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"errors": ["missing client id"]
			}`))
		})

		It("returns a 404 when the campaign type does not exist", func() {
			campaignTypesCollection.GetCall.Returns.Err = collections.NotFoundError{errors.New("campaign type not found")}

			var err error
			request, err = http.NewRequest("GET", "/senders/some-sender-id/campaign_types/missing-campaign-type-id", nil)
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusNotFound))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"errors": ["campaign type not found"]
			}`))
		})

		It("returns a 404 when the sender does not exist", func() {
			campaignTypesCollection.GetCall.Returns.Err = collections.NotFoundError{errors.New("sender not found")}

			var err error
			request, err = http.NewRequest("GET", "/senders/missing-sender-id/campaign_types/some-campaign-type-id", nil)
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusNotFound))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"errors": ["sender not found"]
			}`))
		})

		It("returns a 500 when the collection indicates a system error", func() {
			campaignTypesCollection.GetCall.Returns.Err = errors.New("BOOM!")

			var err error
			request, err = http.NewRequest("GET", "/senders/some-sender-id/campaign_types/some-campaign-type-id", nil)
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusInternalServerError))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"errors": ["BOOM!"]
			}`))
		})
	})
})
