package campaigntypes_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

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

		campaignTypesCollection = fakes.NewCampaignTypesCollection()

		handler = campaigntypes.NewUpdateHandler(campaignTypesCollection)
	})

	It("updates an existing campaign type", func() {
		guid, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())

		campaignTypesCollection.SetCall.ReturnCampaignType = collections.CampaignType{
			ID:          guid.String(),
			Name:        "update-campaign-type",
			Description: "update-campaign-type-description",
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

		request, err = http.NewRequest("PUT", fmt.Sprintf("/senders/some-sender-id/campaign_types/%s", guid.String()), bytes.NewBuffer(requestBody))
		Expect(err).NotTo(HaveOccurred())

		handler.ServeHTTP(writer, request, context)

		Expect(campaignTypesCollection.SetCall.CampaignType).To(Equal(collections.CampaignType{
			ID:          guid.String(),
			Name:        "update-campaign-type",
			Description: "update-campaign-type-description",
			Critical:    true,
			TemplateID:  "",
			SenderID:    "some-sender-id",
		}))

		Expect(writer.Code).To(Equal(http.StatusOK))
		Expect(writer.Body.String()).To(MatchJSON(`{
			"id": "` + guid.String() + `",
			"name": "update-campaign-type",
			"description": "update-campaign-type-description",
			"critical": true,
			"template_id": ""
		}`))
	})
})
