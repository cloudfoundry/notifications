package campaigntypes_test

import (
	"encoding/json"

	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/web/campaigntypes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CampaignTypeResponse", func() {
	It("provides a JSON representation of a campaign type resource", func() {
		campaignType := collections.CampaignType{
			ID:          "some-campaign-type-id",
			Name:        "some-campaign-type",
			Description: "cool campaign type",
			Critical:    true,
			TemplateID:  "some-template-id",
		}

		response := campaigntypes.NewCampaignTypeResponse(campaignType)
		Expect(response).To(Equal(campaigntypes.CampaignTypeResponse{
			ID:          "some-campaign-type-id",
			Name:        "some-campaign-type",
			Description: "cool campaign type",
			Critical:    true,
			TemplateID:  "some-template-id",
			Links: campaigntypes.CampaignTypeResponseLinks{
				Self: campaigntypes.Link{"/campaign_types/some-campaign-type-id"},
			},
		}))
	})

	It("can marshal into JSON", func() {
		campaignType := collections.CampaignType{
			ID:          "some-campaign-type-id",
			Name:        "some-campaign-type",
			Description: "cool campaign type",
			Critical:    true,
			TemplateID:  "some-template-id",
		}

		output, err := json.Marshal(campaigntypes.NewCampaignTypeResponse(campaignType))
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(MatchJSON(`{
			"id":   "some-campaign-type-id",
			"name": "some-campaign-type",
			"description": "cool campaign type",
			"critical": true,
			"template_id": "some-template-id",
			"_links": {
				"self": {
					"href": "/campaign_types/some-campaign-type-id"
				}
			}
		}`))
	})
})
