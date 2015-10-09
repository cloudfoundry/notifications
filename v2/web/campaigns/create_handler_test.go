package campaigns_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"time"

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

var _ = Describe("CreateHandler", func() {
	var (
		handler             campaigns.CreateHandler
		campaignsCollection *mocks.CampaignsCollection
		context             stack.Context
		writer              *httptest.ResponseRecorder
		request             *http.Request
		database            *mocks.Database
		conn                *mocks.Connection
		clock               *mocks.Clock
		startTime           time.Time
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

		startTime = time.Now()
		clock = mocks.NewClock()
		clock.NowCall.Returns.Time = startTime

		context = stack.NewContext()
		context.Set("token", token)
		context.Set("database", database)
		context.Set("client_id", "my-client")

		campaignsCollection = mocks.NewCampaignsCollection()
		campaignsCollection.CreateCall.Returns.Campaign = collections.Campaign{
			ID: "my-campaign-id",
			SendTo: map[string][]string{
				"users": {"user-123", "user-456"},
			},
			CampaignTypeID: "some-campaign-type-id",
			Text:           "come see our new stuff",
			HTML:           "<h1>New stuff</h1>",
			Subject:        "Cool New Stuff",
			TemplateID:     "random-template-id",
			ReplyTo:        "reply-to-address",
		}

		writer = httptest.NewRecorder()

		handler = campaigns.NewCreateHandler(campaignsCollection, clock)
	})

	It("sends a campaign to a list of users", func() {
		requestBody, err := json.Marshal(map[string]interface{}{
			"send_to": map[string][]string{
				"users": {"user-123", "user-456"},
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

		handler.ServeHTTP(writer, request, context)

		Expect(writer.Code).To(Equal(http.StatusAccepted))
		Expect(writer.Body.String()).To(MatchJSON(`{
			"id": "my-campaign-id",
			"send_to": {
				"users": [
					"user-123",
					"user-456"
				]
			},
			"campaign_type_id": "some-campaign-type-id",
			"text":             "come see our new stuff",
			"html":             "<h1>New stuff</h1>",
			"subject":          "Cool New Stuff",
			"template_id":      "random-template-id",
			"reply_to":         "reply-to-address",
			"_links": {
				"self": {"href":"/campaigns/my-campaign-id"},
				"template": {"href":"/templates/random-template-id"},
				"campaign_type": {"href":"/campaign_types/some-campaign-type-id"},
				"status": {"href":"/campaigns/my-campaign-id/status"}
			}
		}`))

		Expect(campaignsCollection.CreateCall.Receives.Connection).To(Equal(conn))
		Expect(campaignsCollection.CreateCall.Receives.ClientID).To(Equal("my-client"))
		Expect(campaignsCollection.CreateCall.Receives.Campaign).To(Equal(collections.Campaign{
			SendTo:         map[string][]string{"users": {"user-123", "user-456"}},
			CampaignTypeID: "some-campaign-type-id",
			Text:           "come see our new stuff",
			HTML:           "<h1>New stuff</h1>",
			Subject:        "Cool New Stuff",
			TemplateID:     "random-template-id",
			ReplyTo:        "reply-to-address",
			SenderID:       "some-sender-id",
			StartTime:      startTime,
		}))
	})

	It("sends a campaign to a list of spaces", func() {
		campaignsCollection.CreateCall.Returns.Campaign.SendTo = map[string][]string{"spaces": {"space-123", "space-456"}}
		requestBody, err := json.Marshal(map[string]interface{}{
			"send_to": map[string][]string{
				"spaces": {"space-123", "space-456"},
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

		handler.ServeHTTP(writer, request, context)

		Expect(writer.Code).To(Equal(http.StatusAccepted))
		Expect(writer.Body.String()).To(MatchJSON(`{
			"id": "my-campaign-id",
			"send_to": {
				"spaces": [
					"space-123",
					"space-456"
				]
			},
			"campaign_type_id": "some-campaign-type-id",
			"text":             "come see our new stuff",
			"html":             "<h1>New stuff</h1>",
			"subject":          "Cool New Stuff",
			"template_id":      "random-template-id",
			"reply_to":         "reply-to-address",
			"_links": {
				"self": {"href":"/campaigns/my-campaign-id"},
				"template": {"href":"/templates/random-template-id"},
				"campaign_type": {"href":"/campaign_types/some-campaign-type-id"},
				"status": {"href":"/campaigns/my-campaign-id/status"}
			}
		}`))

		Expect(campaignsCollection.CreateCall.Receives.Connection).To(Equal(database.Connection()))
		Expect(campaignsCollection.CreateCall.Receives.Campaign).To(Equal(collections.Campaign{
			SendTo:         map[string][]string{"spaces": {"space-123", "space-456"}},
			CampaignTypeID: "some-campaign-type-id",
			Text:           "come see our new stuff",
			HTML:           "<h1>New stuff</h1>",
			Subject:        "Cool New Stuff",
			TemplateID:     "random-template-id",
			ReplyTo:        "reply-to-address",
			SenderID:       "some-sender-id",
			StartTime:      startTime,
		}))
	})

	It("sends a campaign to a list of orgs", func() {
		campaignsCollection.CreateCall.Returns.Campaign.SendTo = map[string][]string{"orgs": {"org-123", "org-456"}}
		requestBody, err := json.Marshal(map[string]interface{}{
			"send_to": map[string][]string{
				"orgs": {"org-123", "org-456"},
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

		handler.ServeHTTP(writer, request, context)

		Expect(writer.Code).To(Equal(http.StatusAccepted))
		Expect(writer.Body.String()).To(MatchJSON(`{
			"id": "my-campaign-id",
			"send_to": {
				"orgs": [
					"org-123",
					"org-456"
				]
			},
			"campaign_type_id": "some-campaign-type-id",
			"text":             "come see our new stuff",
			"html":             "<h1>New stuff</h1>",
			"subject":          "Cool New Stuff",
			"template_id":      "random-template-id",
			"reply_to":         "reply-to-address",
			"_links": {
				"self": {"href":"/campaigns/my-campaign-id"},
				"template": {"href":"/templates/random-template-id"},
				"campaign_type": {"href":"/campaign_types/some-campaign-type-id"},
				"status": {"href":"/campaigns/my-campaign-id/status"}
			}
		}`))

		Expect(campaignsCollection.CreateCall.Receives.Connection).To(Equal(database.Connection()))
		Expect(campaignsCollection.CreateCall.Receives.Campaign).To(Equal(collections.Campaign{
			SendTo:         map[string][]string{"orgs": {"org-123", "org-456"}},
			CampaignTypeID: "some-campaign-type-id",
			Text:           "come see our new stuff",
			HTML:           "<h1>New stuff</h1>",
			Subject:        "Cool New Stuff",
			TemplateID:     "random-template-id",
			ReplyTo:        "reply-to-address",
			SenderID:       "some-sender-id",
			StartTime:      startTime,
		}))
	})

	It("sends a campaign to a list of emails", func() {
		campaignsCollection.CreateCall.Returns.Campaign.SendTo = map[string][]string{"emails": {"test1@example.com", "test2@example.com"}}
		requestBody, err := json.Marshal(map[string]interface{}{
			"send_to": map[string][]string{
				"emails": {"test1@example.com", "test2@example.com"},
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

		handler.ServeHTTP(writer, request, context)

		Expect(writer.Code).To(Equal(http.StatusAccepted))
		Expect(writer.Body.String()).To(MatchJSON(`{
			"id": "my-campaign-id",
			"send_to": {
				"emails": [
					"test1@example.com",
					"test2@example.com"
				]
			},
			"campaign_type_id": "some-campaign-type-id",
			"text":             "come see our new stuff",
			"html":             "<h1>New stuff</h1>",
			"subject":          "Cool New Stuff",
			"template_id":      "random-template-id",
			"reply_to":         "reply-to-address",
			"_links": {
				"self": {"href":"/campaigns/my-campaign-id"},
				"template": {"href":"/templates/random-template-id"},
				"campaign_type": {"href":"/campaign_types/some-campaign-type-id"},
				"status": {"href":"/campaigns/my-campaign-id/status"}
			}
		}`))

		Expect(campaignsCollection.CreateCall.Receives.Connection).To(Equal(database.Connection()))
		Expect(campaignsCollection.CreateCall.Receives.Campaign).To(Equal(collections.Campaign{
			SendTo:         map[string][]string{"emails": {"test1@example.com", "test2@example.com"}},
			CampaignTypeID: "some-campaign-type-id",
			Text:           "come see our new stuff",
			HTML:           "<h1>New stuff</h1>",
			Subject:        "Cool New Stuff",
			TemplateID:     "random-template-id",
			ReplyTo:        "reply-to-address",
			SenderID:       "some-sender-id",
			StartTime:      startTime,
		}))
	})

	Context("when validating user-input", func() {
		Context("when the campaign_type_id is missing", func() {
			BeforeEach(func() {
				requestBody, err := json.Marshal(map[string]interface{}{
					"send_to": map[string][]string{
						"users": {"user-123"},
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
					"send_to": map[string][]string{
						"users": {"user-123"},
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
					"send_to": map[string][]string{
						"users": {"user-123"},
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
					"send_to": map[string][]string{
						"userZ": {"something-obviously-wrong"},
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

		Context("when the email address is invalid", func() {
			BeforeEach(func() {
				requestBody, err := json.Marshal(map[string]interface{}{
					"send_to": map[string][]string{
						"emails": {"malformed-email"},
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
				Expect(writer.Body.String()).To(MatchJSON(`{"errors": ["\"malformed-email\" is not a valid email address"]}`))
			})
		})
	})

	Context("when the token does not have the critical scope", func() {
		BeforeEach(func() {
			requestBody, err := json.Marshal(map[string]interface{}{
				"send_to": map[string][]string{
					"users": {"some-user-guid"},
				},
				"campaign_type_id": "some-campaign-type-id",
				"text":             "come see our new stuff",
				"subject":          "Cool New Stuff",
			})
			Expect(err).NotTo(HaveOccurred())

			campaignsCollection.CreateCall.Returns.Error = collections.PermissionsError{errors.New("Scope critical_notifications.write is required")}

			request, err = http.NewRequest("POST", "/senders/some-sender-id/campaigns", bytes.NewBuffer(requestBody))
			Expect(err).NotTo(HaveOccurred())

		})

		It("returns a 403 forbidden and states the required scope is missing", func() {
			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(403))
			Expect(writer.Body.String()).To(MatchJSON(`{"errors": ["Scope critical_notifications.write is required"]}`))

			Expect(campaignsCollection.CreateCall.Receives.HasCriticalScope).To(BeFalse())
		})
	})

	Context("when the token does have the critical scope", func() {
		BeforeEach(func() {
			requestBody, err := json.Marshal(map[string]interface{}{
				"send_to": map[string][]string{
					"users": {"some-user-guid"},
				},
				"campaign_type_id": "some-campaign-type-id",
				"text":             "come see our new stuff",
				"subject":          "Cool New Stuff",
			})
			Expect(err).NotTo(HaveOccurred())

			request, err = http.NewRequest("POST", "/senders/some-sender-id/campaigns", bytes.NewBuffer(requestBody))
			Expect(err).NotTo(HaveOccurred())
		})

		It("indicates that the requestor has the critical scope", func() {
			tokenHeader := map[string]interface{}{
				"alg": "FAST",
			}
			tokenClaims := map[string]interface{}{
				"client_id": "some-uaa-client-id",
				"exp":       int64(3404281214),
				"scope":     []string{"critical_notifications.write"},
			}
			token, err := jwt.Parse(helpers.BuildToken(tokenHeader, tokenClaims), func(*jwt.Token) (interface{}, error) {
				return []byte(application.UAAPublicKey), nil
			})
			Expect(err).NotTo(HaveOccurred())
			context.Set("token", token)

			handler.ServeHTTP(writer, request, context)

			Expect(campaignsCollection.CreateCall.Receives.HasCriticalScope).To(BeTrue())
		})
	})

	Context("when an error occurs", func() {
		BeforeEach(func() {
			requestBody, err := json.Marshal(map[string]interface{}{
				"send_to": map[string][]string{
					"users": {"user-123"},
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
		})

		Context("when the collection returns an unknown error", func() {
			It("returns a 500 and the corresponding error", func() {
				campaignsCollection.CreateCall.Returns.Error = errors.New("some fantastic error")
				handler.ServeHTTP(writer, request, context)
				Expect(writer.Code).To(Equal(http.StatusInternalServerError))
				Expect(writer.Body.String()).To(MatchJSON(`{"errors": ["some fantastic error"]}`))
			})
		})

		Context("when the collection returns a not found error", func() {
			It("returns a 404 and the corresponding error", func() {
				campaignsCollection.CreateCall.Returns.Error = collections.NotFoundError{errors.New("the entire datacenter has gone away")}
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
