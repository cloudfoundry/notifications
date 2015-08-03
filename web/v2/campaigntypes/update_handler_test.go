package campaigntypes_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/web/v2/campaigntypes"
	"github.com/dgrijalva/jwt-go"
	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ryanmoran/stack"
)

var _ = Describe("UpdateHandler", func() {
	var (
		handler                 campaigntypes.UpdateHandler
		campaignTypesCollection *fakes.CampaignTypesCollection
		context                 stack.Context
		writer                  *httptest.ResponseRecorder
		request                 *http.Request
		database                *fakes.Database
		tokenHeader             map[string]interface{}
		tokenClaims             map[string]interface{}
		guid                    string
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
			"scope": []string{
				"notifications.write",
				"critical_notifications.write",
			},
		}
		rawToken := fakes.BuildToken(tokenHeader, tokenClaims)
		token, err := jwt.Parse(rawToken, func(*jwt.Token) (interface{}, error) {
			return []byte(application.UAAPublicKey), nil
		})
		Expect(err).NotTo(HaveOccurred())
		context.Set("token", token)

		writer = httptest.NewRecorder()

		g, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		guid = g.String()

		campaignTypesCollection = fakes.NewCampaignTypesCollection()

		campaignTypesCollection.GetCall.ReturnCampaignType = collections.CampaignType{
			ID:          guid,
			Name:        "my old name",
			Description: "old description",
			Critical:    true,
			TemplateID:  "",
			SenderID:    "some-sender-id",
		}

		requestBody, err := json.Marshal(map[string]interface{}{
			"name":        "update-campaign-type",
			"description": "update-campaign-type-description",
			"critical":    true,
			"template_id": "",
		})
		Expect(err).NotTo(HaveOccurred())

		request, err = http.NewRequest("PUT", fmt.Sprintf("/senders/some-sender-id/campaign_types/%s", guid), bytes.NewBuffer(requestBody))
		Expect(err).NotTo(HaveOccurred())

		handler = campaigntypes.NewUpdateHandler(campaignTypesCollection)
	})

	It("updates an existing campaign type", func() {
		campaignTypesCollection.SetCall.ReturnCampaignType = collections.CampaignType{
			ID:          guid,
			Name:        "update-campaign-type",
			Description: "update-campaign-type-description",
			Critical:    true,
			TemplateID:  "",
			SenderID:    "some-sender-id",
		}

		handler.ServeHTTP(writer, request, context)

		Expect(campaignTypesCollection.SetCall.CampaignType).To(Equal(collections.CampaignType{
			ID:          guid,
			Name:        "update-campaign-type",
			Description: "update-campaign-type-description",
			Critical:    true,
			TemplateID:  "",
			SenderID:    "some-sender-id",
		}))

		Expect(writer.Code).To(Equal(http.StatusOK))
		Expect(writer.Body.String()).To(MatchJSON(`{
			"id": "` + guid + `",
			"name": "update-campaign-type",
			"description": "update-campaign-type-description",
			"critical": true,
			"template_id": ""
		}`))
	})

	It("works when only the name field is updated", func() {
		requestBody, err := json.Marshal(map[string]interface{}{
			"name": "my new name",
		})
		Expect(err).NotTo(HaveOccurred())

		campaignTypesCollection.SetCall.ReturnCampaignType = collections.CampaignType{
			ID:          guid,
			Name:        "my new name",
			Description: "old description",
			Critical:    true,
			TemplateID:  "",
			SenderID:    "some-sender-id",
		}

		request, err := http.NewRequest("PUT", fmt.Sprintf("/senders/some-sender-id/campaign_types/%s", guid), bytes.NewBuffer(requestBody))
		Expect(err).NotTo(HaveOccurred())

		handler.ServeHTTP(writer, request, context)

		Expect(campaignTypesCollection.SetCall.CampaignType).To(Equal(collections.CampaignType{
			ID:          guid,
			Name:        "my new name",
			Description: "old description",
			Critical:    true,
			TemplateID:  "",
			SenderID:    "some-sender-id",
		}))

		Expect(writer.Code).To(Equal(http.StatusOK))
		Expect(writer.Body.String()).To(MatchJSON(`{
			"id": "` + guid + `",
			"name": "my new name",
			"description": "old description",
			"critical": true,
			"template_id": ""
		}`))
	})

	It("works when no parameters are passed into the update", func() {
		requestBody, err := json.Marshal(map[string]interface{}{})
		Expect(err).NotTo(HaveOccurred())

		campaignTypesCollection.SetCall.ReturnCampaignType = collections.CampaignType{
			ID:          guid,
			Name:        "my old name",
			Description: "old description",
			Critical:    true,
			TemplateID:  "",
			SenderID:    "some-sender-id",
		}

		request, err := http.NewRequest("PUT", fmt.Sprintf("/senders/some-sender-id/campaign_types/%s", guid), bytes.NewBuffer(requestBody))
		Expect(err).NotTo(HaveOccurred())

		handler.ServeHTTP(writer, request, context)

		Expect(campaignTypesCollection.SetCall.CampaignType).To(Equal(collections.CampaignType{
			ID:          guid,
			Name:        "my old name",
			Description: "old description",
			Critical:    true,
			TemplateID:  "",
			SenderID:    "some-sender-id",
		}))

		Expect(writer.Code).To(Equal(http.StatusOK))
		Expect(writer.Body.String()).To(MatchJSON(`{
			"id": "` + guid + `",
			"name": "my old name",
			"description": "old description",
			"critical": true,
			"template_id": ""
		}`))
	})

	Context("failure cases", func() {
		It("returns a 400 when the request JSON cannot be unmarshalled", func() {
			request.Body = ioutil.NopCloser(strings.NewReader("%%%%"))

			handler.ServeHTTP(writer, request, context)

			Expect(writer.Code).To(Equal(http.StatusBadRequest))
			Expect(writer.Body.String()).To(MatchJSON(`{
					"error": "invalid json body"
			}`))
		})

		It("returns a 422 if the name field is updated to an empty string", func() {
			requestBody, err := json.Marshal(map[string]interface{}{
				"name": "",
			})
			Expect(err).NotTo(HaveOccurred())

			request, err = http.NewRequest("PUT", fmt.Sprintf("/senders/some-sender-id/campaign_types/%s", guid), bytes.NewBuffer(requestBody))

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(422))
			Expect(writer.Body.String()).To(MatchJSON(`{
					"error": "name field cannot be blank"
			}`))
		})

		It("returns a 422 if the description field is updated to an empty string", func() {
			requestBody, err := json.Marshal(map[string]interface{}{
				"description": "",
			})
			Expect(err).NotTo(HaveOccurred())

			request, err = http.NewRequest("PUT", fmt.Sprintf("/senders/some-sender-id/campaign_types/%s", guid), bytes.NewBuffer(requestBody))

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(422))
			Expect(writer.Body.String()).To(MatchJSON(`{
					"error": "description field cannot be blank"
			}`))
		})

		It("returns a 422 if the name and description fields are updated to an empty string", func() {
			requestBody, err := json.Marshal(map[string]interface{}{
				"name":        "",
				"description": "",
			})
			Expect(err).NotTo(HaveOccurred())

			request, err = http.NewRequest("PUT", fmt.Sprintf("/senders/some-sender-id/campaign_types/%s", guid), bytes.NewBuffer(requestBody))

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(422))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"error": "name field cannot be blank, description field cannot be blank"
			}`))
		})

		PIt("returns a 404 when the sender could not be found", func() {
			campaignTypesCollection.SetCall.Err = collections.NotFoundError{
				Err:     errors.New("THIS WAS PRODUCED BY ROBOTS"),
				Message: "This is for humans.",
			}

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusNotFound))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"error": "This is for humans."
			}`))
		})

		It("returns a 404 when the campaign type could not be found", func() {
			campaignTypesCollection.GetCall.Err = collections.NotFoundError{
				Err:     errors.New("THIS WAS PRODUCED BY ROBOTS"),
				Message: "This is for humans.",
			}

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusNotFound))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"error": "This is for humans."
			}`))
		})

		// XIt("does the right thing when critical flag is set true but sender does not have UAA critical-write")

		// XIt("allows an update of Critical from true -> false even if the user does not have UAA critical-write")
	})
})
