package campaigns_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/cloudfoundry-incubator/notifications/testing/helpers"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/web/campaigns"
	"github.com/dgrijalva/jwt-go"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetHandler", func() {
	var (
		handler             campaigns.GetHandler
		campaignsCollection *mocks.CampaignsCollection
		context             stack.Context
		writer              *httptest.ResponseRecorder
		request             *http.Request
		database            *mocks.Database
		conn                *mocks.Connection
	)

	BeforeEach(func() {
		tokenHeader := map[string]interface{}{
			"alg": "FAST",
		}
		tokenClaims := map[string]interface{}{
			"client_id": "some-uaa-client-id",
			"exp":       int64(3404281214),
			"scope":     []string{"notifications.write"},
		}
		token, err := jwt.Parse(helpers.BuildToken(tokenHeader, tokenClaims), func(*jwt.Token) (interface{}, error) {
			return []byte(application.UAAPublicKey), nil
		})
		Expect(err).NotTo(HaveOccurred())

		conn = mocks.NewConnection()
		database = mocks.NewDatabase()
		database.ConnectionCall.Returns.Connection = conn

		context = stack.NewContext()
		context.Set("token", token)
		context.Set("database", database)
		context.Set("client_id", "my-client")

		campaignsCollection = mocks.NewCampaignsCollection()
		campaignsCollection.GetCall.Returns.Campaign = collections.Campaign{
			ID:             "some-campaign-id",
			SendTo:         map[string]string{"user": "user-123"},
			CampaignTypeID: "some-campaign-type-id",
			Text:           "come see our new stuff",
			HTML:           "<h1>New stuff</h1>",
			Subject:        "Cool New Stuff",
			TemplateID:     "random-template-id",
			ReplyTo:        "reply-to-address",
			ClientID:       "my-client",
		}

		writer = httptest.NewRecorder()

		handler = campaigns.NewGetHandler(campaignsCollection)
	})

	It("gets the details about an existing campaign", func() {
		var err error
		request, err = http.NewRequest("GET", "/senders/some-sender-id/campaigns/some-campaign-id", nil)
		Expect(err).NotTo(HaveOccurred())

		handler.ServeHTTP(writer, request, context)

		Expect(writer.Code).To(Equal(http.StatusOK))
		Expect(writer.Body).To(MatchJSON(`{
			"id":               "some-campaign-id",
			"send_to": {
				"user": "user-123"
			},
			"campaign_type_id": "some-campaign-type-id",
			"text":             "come see our new stuff",
			"html":             "<h1>New stuff</h1>",
			"subject":          "Cool New Stuff",
			"template_id":      "random-template-id",
			"reply_to":         "reply-to-address"
		}`))

		Expect(campaignsCollection.GetCall.Receives.Connection).To(Equal(conn))
		Expect(campaignsCollection.GetCall.Receives.ClientID).To(Equal("my-client"))
		Expect(campaignsCollection.GetCall.Receives.SenderID).To(Equal("some-sender-id"))
		Expect(campaignsCollection.GetCall.Receives.CampaignID).To(Equal("some-campaign-id"))
	})

	Context("failure cases", func() {
		It("returns a 404 if the campaign could not be found", func() {
			campaignsCollection.GetCall.Returns.Error = collections.NotFoundError{errors.New("campaign not found")}

			var err error
			request, err = http.NewRequest("GET", "/senders/some-sender-id/campaigns/missing-campaign-id", nil)
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)

			Expect(writer.Code).To(Equal(http.StatusNotFound))
			Expect(writer.Body).To(MatchJSON(`{
				"errors": [
					"campaign not found"
				]
			}`))
		})

		It("returns a 500 if an unknown error occurs", func() {
			campaignsCollection.GetCall.Returns.Error = errors.New("something went wrong")

			var err error
			request, err = http.NewRequest("GET", "/senders/some-sender-id/campaigns/missing-campaign-id", nil)
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)

			Expect(writer.Code).To(Equal(http.StatusInternalServerError))
			Expect(writer.Body).To(MatchJSON(`{
				"errors": [
					"something went wrong"
				]
			}`))
		})
	})
})
