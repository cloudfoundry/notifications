package campaigns_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/web/campaigns"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CreateHandler", func() {
	var (
		handler             campaigns.CreateHandler
		campaignsCollection *mocks.CampaignsCollection
		context             stack.Context
		writer              *httptest.ResponseRecorder
		request             *http.Request
		database            *mocks.Database
	)

	BeforeEach(func() {
		context = stack.NewContext()

		database = mocks.NewDatabase()
		context.Set("database", database)

		context.Set("client_id", "my-client")

		campaignsCollection = mocks.NewCampaignsCollection()
		campaignsCollection.CreateCall.Returns.Campaign = collections.Campaign{
			ID: "my-campaign-id",
		}

		writer = httptest.NewRecorder()

		requestBody, err := json.Marshal(map[string]interface{}{
			"send_to": map[string]string{
				"user": "user-123",
			},
			"campaign_type_id": "some-campaign-type-id",
			"text":             "come see our new stuff",
			"html":             "<h1>New stuff</h1>",
			"subject":          "Cool New Stuff",
			"template_id":      "random-template-id",
			"reply_to":         "reply-to-address",
		})
		Expect(err).NotTo(HaveOccurred())

		request, err = http.NewRequest("POST", "/senders/some-sender-id/campaigns", bytes.NewBuffer(requestBody))
		Expect(err).NotTo(HaveOccurred())

		handler = campaigns.NewCreateHandler(campaignsCollection)
	})

	It("sends a campaign to a user", func() {
		handler.ServeHTTP(writer, request, context)

		Expect(campaignsCollection.CreateCall.Receives.Conn).To(Equal(database.Connection()))
		Expect(campaignsCollection.CreateCall.Receives.Campaign).To(Equal(collections.Campaign{
			SendTo:         map[string]string{"user": "user-123"},
			CampaignTypeID: "some-campaign-type-id",
			Text:           "come see our new stuff",
			HTML:           "<h1>New stuff</h1>",
			Subject:        "Cool New Stuff",
			TemplateID:     "random-template-id",
			ReplyTo:        "reply-to-address",
			ClientID:       "my-client",
		}))
		Expect(writer.Code).To(Equal(http.StatusAccepted))
		Expect(writer.Body.String()).To(MatchJSON(`{
			"campaign_id": "my-campaign-id"
		}`))
	})

	Context("when validating user-input", func() {
		Context("when the campaign_type_id is missing", func() {
			BeforeEach(func() {
				requestBody, err := json.Marshal(map[string]interface{}{
					"send_to": map[string]string{
						"user": "user-123",
					},
					"campaign_type_id": "",
					"text":             "come see our new stuff",
					"subject":          "Cool New Stuff",
					"template_id":      "random-template-id",
					"reply_to":         "reply-to-address",
				})
				Expect(err).NotTo(HaveOccurred())

				request, err = http.NewRequest("POST", "/senders/some-sender-id/campaigns", bytes.NewBuffer(requestBody))
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns a 422 and states that the request is missing a campaign type", func() {
				handler.ServeHTTP(writer, request, context)
				Expect(writer.Code).To(Equal(422))
				Expect(writer.Body.String()).To(MatchJSON(`{"errors": ["missing campaign_type_id"]}`))
			})
		})

		Context("when both the text and html bodies are missing", func() {
			BeforeEach(func() {
				requestBody, err := json.Marshal(map[string]interface{}{
					"send_to": map[string]string{
						"user": "user-123",
					},
					"campaign_type_id": "some-campaign-type-id",
					"text":             "",
					"html":             "",
					"subject":          "Cool New Stuff",
					"template_id":      "random-template-id",
					"reply_to":         "reply-to-address",
				})
				Expect(err).NotTo(HaveOccurred())

				request, err = http.NewRequest("POST", "/senders/some-sender-id/campaigns", bytes.NewBuffer(requestBody))
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns a 422 and states that the request is missing either text or html", func() {
				handler.ServeHTTP(writer, request, context)
				Expect(writer.Code).To(Equal(422))
				Expect(writer.Body.String()).To(MatchJSON(`{"errors": ["missing either campaign text or html"]}`))
			})
		})

		Context("when the subject is missing", func() {
			BeforeEach(func() {
				requestBody, err := json.Marshal(map[string]interface{}{
					"send_to": map[string]string{
						"user": "user-123",
					},
					"campaign_type_id": "some-campaign-type-id",
					"text":             "come see our new stuff",
					"subject":          "",
					"template_id":      "random-template-id",
					"reply_to":         "reply-to-address",
				})
				Expect(err).NotTo(HaveOccurred())

				request, err = http.NewRequest("POST", "/senders/some-sender-id/campaigns", bytes.NewBuffer(requestBody))
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns a 422 and states that the request is missing a subject", func() {
				handler.ServeHTTP(writer, request, context)
				Expect(writer.Code).To(Equal(422))
				Expect(writer.Body.String()).To(MatchJSON(`{"errors": ["missing subject"]}`))
			})
		})

		Context("when the audience specific key is invalid", func() {
			BeforeEach(func() {
				requestBody, err := json.Marshal(map[string]interface{}{
					"send_to": map[string]string{
						"userZ": "something-obviously-wrong",
					},
					"campaign_type_id": "some-campaign-type-id",
					"text":             "come see our new stuff",
					"subject":          "Cool New Stuff",
				})
				Expect(err).NotTo(HaveOccurred())

				request, err = http.NewRequest("POST", "/senders/some-sender-id/campaigns", bytes.NewBuffer(requestBody))
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns a 422 and states the audience is invalid", func() {
				handler.ServeHTTP(writer, request, context)
				Expect(writer.Code).To(Equal(422))
				Expect(writer.Body.String()).To(MatchJSON(`{"errors": ["\"userZ\" is not a valid audience"]}`))
			})
		})
	})

	Context("when an error occurs", func() {
		Context("when the collection returns an unknown error", func() {
			It("returns a 500 and the corresponding error", func() {
				campaignsCollection.CreateCall.Returns.Err = errors.New("some fantastic error")
				handler.ServeHTTP(writer, request, context)
				Expect(writer.Code).To(Equal(http.StatusInternalServerError))
				Expect(writer.Body.String()).To(MatchJSON(`{"errors": ["some fantastic error"]}`))
			})
		})

		Context("when the collection returns a not found error", func() {
			It("returns a 404 and the corresponding error", func() {
				campaignsCollection.CreateCall.Returns.Err = collections.NotFoundError{errors.New("the entire datacenter has gone away")}
				handler.ServeHTTP(writer, request, context)
				Expect(writer.Code).To(Equal(http.StatusNotFound))
				Expect(writer.Body.String()).To(MatchJSON(`{"errors": ["the entire datacenter has gone away"]}`))
			})
		})

		Context("when the request JSON is not well-formed", func() {
			It("returns a 400 and states that the request is invalid", func() {
				request, err := http.NewRequest("POST", "/senders/some-sender-id/campaigns", bytes.NewBufferString("%%%"))
				Expect(err).NotTo(HaveOccurred())

				handler.ServeHTTP(writer, request, context)
				Expect(writer.Code).To(Equal(http.StatusBadRequest))
				Expect(writer.Body.String()).To(MatchJSON(`{"errors": ["invalid json body"]}`))
			})
		})
	})
})
