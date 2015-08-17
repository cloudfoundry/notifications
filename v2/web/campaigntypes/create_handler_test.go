package campaigntypes_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/cloudfoundry-incubator/notifications/testing/fakes"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/web/campaigntypes"
	"github.com/dgrijalva/jwt-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ryanmoran/stack"
)

var _ = Describe("CreateHandler", func() {
	var (
		handler                 campaigntypes.CreateHandler
		campaignTypesCollection *fakes.CampaignTypesCollection
		context                 stack.Context
		writer                  *httptest.ResponseRecorder
		request                 *http.Request
		database                *fakes.Database
		tokenHeader             map[string]interface{}
		tokenClaims             map[string]interface{}
	)

	BeforeEach(func() {
		context = stack.NewContext()

		context.Set("client_id", "some-client-id")

		database = fakes.NewDatabase()
		context.Set("database", database)

		tokenHeader = map[string]interface{}{
			"alg": "FAST",
		}
		tokenClaims = map[string]interface{}{
			"client_id": "some-uaa-client-id",
			"exp":       int64(3404281214),
			"scope":     []string{"notifications.write"},
		}
		rawToken := fakes.BuildToken(tokenHeader, tokenClaims)
		token, err := jwt.Parse(rawToken, func(*jwt.Token) (interface{}, error) {
			return []byte(application.UAAPublicKey), nil
		})
		Expect(err).NotTo(HaveOccurred())
		context.Set("token", token)

		writer = httptest.NewRecorder()
		campaignTypesCollection = fakes.NewCampaignTypesCollection()
		campaignTypesCollection.SetCall.Returns.CampaignType = collections.CampaignType{
			ID:          "some-campaign-type-id",
			Name:        "some-campaign-type",
			Description: "some-campaign-type-description",
			Critical:    false,
			TemplateID:  "some-template-id",
			SenderID:    "some-sender-id",
		}

		requestBody, err := json.Marshal(map[string]interface{}{
			"name":        "some-campaign-type",
			"description": "some-campaign-type-description",
			"critical":    false,
			"template_id": "some-template-id",
		})
		Expect(err).NotTo(HaveOccurred())

		request, err = http.NewRequest("POST", "/senders/some-sender-id/campaign_types", bytes.NewBuffer(requestBody))
		Expect(err).NotTo(HaveOccurred())

		handler = campaigntypes.NewCreateHandler(campaignTypesCollection)
	})

	It("creates a campaign type", func() {
		handler.ServeHTTP(writer, request, context)

		Expect(campaignTypesCollection.SetCall.Receives.CampaignType).To(Equal(collections.CampaignType{
			Name:        "some-campaign-type",
			Description: "some-campaign-type-description",
			Critical:    false,
			TemplateID:  "some-template-id",
			SenderID:    "some-sender-id",
		}))
		Expect(campaignTypesCollection.SetCall.Receives.Conn).To(Equal(database.Conn))
		Expect(database.ConnectionWasCalled).To(BeTrue())

		Expect(writer.Code).To(Equal(http.StatusCreated))
		Expect(writer.Body.String()).To(MatchJSON(`{
			"id": "some-campaign-type-id",
			"name": "some-campaign-type",
			"description": "some-campaign-type-description",
			"critical": false,
			"template_id": "some-template-id"
		}`))
	})

	It("requires critical_notifications.write to create a critical campaign type", func() {
		tokenClaims["scope"] = []string{"notifications.write", "critical_notifications.write"}
		rawToken := fakes.BuildToken(tokenHeader, tokenClaims)
		token, err := jwt.Parse(rawToken, func(*jwt.Token) (interface{}, error) {
			return []byte(application.UAAPublicKey), nil
		})
		Expect(err).NotTo(HaveOccurred())
		context.Set("token", token)

		campaignTypesCollection.SetCall.Returns.CampaignType = collections.CampaignType{
			ID:          "some-campaign-type-id",
			Name:        "some-campaign-type",
			Description: "some-campaign-type-description",
			Critical:    true,
			TemplateID:  "some-template-id",
			SenderID:    "some-sender-id",
		}

		requestBody, err := json.Marshal(map[string]interface{}{
			"name":        "some-campaign-type",
			"description": "some-campaign-type-description",
			"critical":    true,
			"template_id": "some-template-id",
		})
		Expect(err).NotTo(HaveOccurred())

		request, err = http.NewRequest("POST", "/senders/some-sender-id/campaign_types", bytes.NewBuffer(requestBody))
		Expect(err).NotTo(HaveOccurred())

		handler.ServeHTTP(writer, request, context)

		Expect(campaignTypesCollection.SetCall.Receives.CampaignType).To(Equal(collections.CampaignType{
			Name:        "some-campaign-type",
			Description: "some-campaign-type-description",
			Critical:    true,
			TemplateID:  "some-template-id",
			SenderID:    "some-sender-id",
		}))

		Expect(writer.Code).To(Equal(http.StatusCreated))
		Expect(writer.Body.String()).To(MatchJSON(`{
			"id": "some-campaign-type-id",
			"name": "some-campaign-type",
			"description": "some-campaign-type-description",
			"critical": true,
			"template_id": "some-template-id"
		}`))
	})

	Context("failure cases", func() {
		It("returns a 403 when the client without the critical_notifications.write scope attempts to create a critical campaign type", func() {
			campaignTypesCollection.SetCall.Returns.CampaignType = collections.CampaignType{
				ID:          "some-campaign-type-id",
				Name:        "some-campaign-type",
				Description: "some-campaign-type-description",
				Critical:    true,
				TemplateID:  "some-template-id",
				SenderID:    "some-sender-id",
			}

			requestBody, err := json.Marshal(map[string]interface{}{
				"name":        "some-campaign-type",
				"description": "some-campaign-type-description",
				"critical":    true,
				"template_id": "some-template-id",
			})
			Expect(err).NotTo(HaveOccurred())

			request, err = http.NewRequest("POST", "/senders/some-sender-id/campaign_types", bytes.NewBuffer(requestBody))
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusForbidden))
			Expect(writer.Body.String()).To(MatchJSON(`{ "errors": ["You do not have permission to create critical campaign types"] }`))
		})

		It("returns a 400 when the JSON request body cannot be unmarshalled", func() {
			var err error
			request, err = http.NewRequest("POST", "/senders/some-sender-id/campaign_types", strings.NewReader("%%%"))
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusBadRequest))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"errors": ["invalid json body"]
			}`))
		})

		It("returns a 422 when name is omitted", func() {
			requestBody, err := json.Marshal(map[string]interface{}{
				"description": "description",
				"critical":    false,
				"template_id": "",
			})
			Expect(err).NotTo(HaveOccurred())

			request, err = http.NewRequest("POST", "/senders/some-sender-id/campaign_types", bytes.NewBuffer(requestBody))
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(422))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"errors": ["missing campaign type name"]
			}`))
		})

		It("returns a 422 when description is omitted", func() {
			requestBody, err := json.Marshal(map[string]interface{}{
				"name":        "some name",
				"critical":    false,
				"template_id": "",
			})
			Expect(err).NotTo(HaveOccurred())

			request, err = http.NewRequest("POST", "/senders/some-sender-id/campaign_types", bytes.NewBuffer(requestBody))
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(422))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"errors": ["missing campaign type description"]
			}`))
		})

		It("returns a 500 when there is a persistence error", func() {
			campaignTypesCollection.SetCall.Returns.Err = errors.New("BOOM!")

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusInternalServerError))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"errors": ["BOOM!"]
			}`))
		})

		It("returns a 404 when the collection returns a NotFoundError", func() {
			campaignTypesCollection.SetCall.Returns.Err = collections.NotFoundError{errors.New("something not found")}

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusNotFound))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"errors": ["something not found"]
			}`))
		})
	})
})
