package campaigntypes_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/helpers"
	"github.com/cloudfoundry-incubator/notifications/web/v2/campaigntypes"
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
		campaignTypesCollection.AddCall.ReturnCampaignType = collections.CampaignType{
			ID:          "some-campaign-type-id",
			Name:        helpers.AddressOfString("some-campaign-type"),
			Description: helpers.AddressOfString("some-campaign-type-description"),
			Critical:    helpers.AddressOfBool(false),
			TemplateID:  helpers.AddressOfString("some-template-id"),
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

		Expect(campaignTypesCollection.AddCall.CampaignType).To(Equal(collections.CampaignType{
			Name:        helpers.AddressOfString("some-campaign-type"),
			Description: helpers.AddressOfString("some-campaign-type-description"),
			Critical:    helpers.AddressOfBool(false),
			TemplateID:  helpers.AddressOfString("some-template-id"),
			SenderID:    "some-sender-id",
		}))
		Expect(campaignTypesCollection.AddCall.Conn).To(Equal(database.Conn))
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

	It("allows the critical field to be omitted", func() {
		requestBody, err := json.Marshal(map[string]string{
			"name":        "some-campaign-type",
			"description": "some-campaign-type-description",
			"template_id": "some-template-id",
		})
		Expect(err).NotTo(HaveOccurred())

		request, err = http.NewRequest("POST", "/senders/some-sender-id/campaign_types", bytes.NewBuffer(requestBody))
		Expect(err).NotTo(HaveOccurred())

		handler.ServeHTTP(writer, request, context)

		Expect(campaignTypesCollection.AddCall.CampaignType).To(Equal(collections.CampaignType{
			Name:        helpers.AddressOfString("some-campaign-type"),
			Description: helpers.AddressOfString("some-campaign-type-description"),
			Critical:    nil,
			TemplateID:  helpers.AddressOfString("some-template-id"),
			SenderID:    "some-sender-id",
		}))

		Expect(writer.Code).To(Equal(http.StatusCreated))
		Expect(writer.Body.String()).To(MatchJSON(`{
			"id": "some-campaign-type-id",
			"name": "some-campaign-type",
			"description": "some-campaign-type-description",
			"critical": false,
			"template_id": "some-template-id"
		}`))
	})

	It("allows the template_id field to be omitted", func() {
		campaignTypesCollection.AddCall.ReturnCampaignType = collections.CampaignType{
			ID:          "some-campaign-type-id",
			Name:        helpers.AddressOfString("some-campaign-type"),
			Description: helpers.AddressOfString("some-campaign-type-description"),
			Critical:    helpers.AddressOfBool(false),
			TemplateID:  helpers.AddressOfString(""),
			SenderID:    "some-sender-id",
		}

		requestBody, err := json.Marshal(map[string]interface{}{
			"name":        "some-campaign-type",
			"description": "some-campaign-type-description",
			"critical":    false,
		})
		Expect(err).NotTo(HaveOccurred())

		request, err = http.NewRequest("POST", "/senders/some-sender-id/campaign_types", bytes.NewBuffer(requestBody))
		Expect(err).NotTo(HaveOccurred())

		handler.ServeHTTP(writer, request, context)

		Expect(campaignTypesCollection.AddCall.CampaignType).To(Equal(collections.CampaignType{
			Name:        helpers.AddressOfString("some-campaign-type"),
			Description: helpers.AddressOfString("some-campaign-type-description"),
			Critical:    helpers.AddressOfBool(false),
			TemplateID:  nil,
			SenderID:    "some-sender-id",
		}))

		Expect(writer.Code).To(Equal(http.StatusCreated))
		Expect(writer.Body.String()).To(MatchJSON(`{
			"id": "some-campaign-type-id",
			"name": "some-campaign-type",
			"description": "some-campaign-type-description",
			"critical": false,
			"template_id": ""
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

		campaignTypesCollection.AddCall.ReturnCampaignType = collections.CampaignType{
			ID:          "some-campaign-type-id",
			Name:        helpers.AddressOfString("some-campaign-type"),
			Description: helpers.AddressOfString("some-campaign-type-description"),
			Critical:    helpers.AddressOfBool(true),
			TemplateID:  helpers.AddressOfString("some-template-id"),
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

		Expect(campaignTypesCollection.AddCall.CampaignType).To(Equal(collections.CampaignType{
			Name:        helpers.AddressOfString("some-campaign-type"),
			Description: helpers.AddressOfString("some-campaign-type-description"),
			Critical:    helpers.AddressOfBool(true),
			TemplateID:  helpers.AddressOfString("some-template-id"),
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
			campaignTypesCollection.AddCall.ReturnCampaignType = collections.CampaignType{
				ID:          "some-campaign-type-id",
				Name:        helpers.AddressOfString("some-campaign-type"),
				Description: helpers.AddressOfString("some-campaign-type-description"),
				Critical:    helpers.AddressOfBool(true),
				TemplateID:  helpers.AddressOfString("some-template-id"),
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
			Expect(writer.Body.String()).To(MatchJSON(`{ "error": "Forbidden" }`))
		})

		It("returns a 400 when the JSON request body cannot be unmarshalled", func() {
			var err error
			request, err = http.NewRequest("POST", "/senders/some-sender-id/campaign_types", strings.NewReader("%%%"))
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusBadRequest))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"error": "invalid json body"
			}`))
		})

		It("returns a 422 when the model does not save", func() {
			requestBody, err := json.Marshal(map[string]interface{}{
				"description": "missing-name",
				"critical":    false,
				"template_id": "some-template-id",
			})
			Expect(err).NotTo(HaveOccurred())

			request, err = http.NewRequest("POST", "/senders/some-sender-id/campaign_types", bytes.NewBuffer(requestBody))
			Expect(err).NotTo(HaveOccurred())

			campaignTypesCollection.AddCall.Err = collections.NewValidationError("bananas are delicious")

			handler.ServeHTTP(writer, request, context)

			Expect(writer.Code).To(Equal(422))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"error": "bananas are delicious"
			}`))
		})

		It("returns a 500 when there is a persistence error", func() {
			campaignTypesCollection.AddCall.Err = errors.New("BOOM!")

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusInternalServerError))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"error": "BOOM!"
			}`))
		})
		It("returns a 404 when the sender could not be found", func() {
			campaignTypesCollection.AddCall.Err = collections.NotFoundError{
				Err:     errors.New("THIS WAS PRODUCED BY ROBOTS"),
				Message: "This is for humans.",
			}

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusNotFound))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"error": "This is for humans."
			}`))
		})

		PIt("returns a 422 when the template_id is not valid for the given client", func() {})
	})
})
