package campaigntypes_test

import (
	"encoding/json"

	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/web/campaigntypes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CampaignTypesListResponse", func() {
	var response campaigntypes.CampaignTypesListResponse

	BeforeEach(func() {
		response = campaigntypes.NewCampaignTypesListResponse("some-sender-id", []collections.CampaignType{
			{
				ID:          "some-campaign-type-id",
				Name:        "some-campaign-type",
				Description: "first campaign type",
				Critical:    false,
				TemplateID:  "",
			},
			{
				ID:          "another-campaign-type-id",
				Name:        "another-campaign-type",
				Description: "second campaign type",
				Critical:    true,
				TemplateID:  "template-id",
			},
		})
	})

	It("provides a JSON representation of a list of campaign_type resources", func() {
		Expect(response).To(Equal(campaigntypes.CampaignTypesListResponse{
			CampaignTypes: []campaigntypes.CampaignTypeResponse{
				{
					ID:          "some-campaign-type-id",
					Name:        "some-campaign-type",
					Description: "first campaign type",
					Critical:    false,
					TemplateID:  "",
					Links: campaigntypes.CampaignTypeResponseLinks{
						Self: campaigntypes.Link{"/campaign_types/some-campaign-type-id"},
					},
				},
				{
					ID:          "another-campaign-type-id",
					Name:        "another-campaign-type",
					Description: "second campaign type",
					Critical:    true,
					TemplateID:  "template-id",
					Links: campaigntypes.CampaignTypeResponseLinks{
						Self: campaigntypes.Link{"/campaign_types/another-campaign-type-id"},
					},
				},
			},
			Links: campaigntypes.CampaignTypesListResponseLinks{
				Self:   campaigntypes.Link{"/senders/some-sender-id/campaign_types"},
				Sender: campaigntypes.Link{"/senders/some-sender-id"},
			},
		}))
	})

	It("can marshal into JSON", func() {
		output, err := json.Marshal(response)
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(MatchJSON(`{
			"campaign_types": [
				{
					"id":   "some-campaign-type-id",
					"name": "some-campaign-type",
					"description": "first campaign type",
					"critical": false,
					"template_id": "",
					"_links": {
						"self": {
							"href": "/campaign_types/some-campaign-type-id"
						}
					}
				},
				{
					"id":   "another-campaign-type-id",
					"name": "another-campaign-type",
					"description": "second campaign type",
					"critical": true,
					"template_id": "template-id",
					"_links": {
						"self": {
							"href": "/campaign_types/another-campaign-type-id"
						}
					}
				}
			],
			"_links": {
				"self": {
					"href": "/senders/some-sender-id/campaign_types"
				},
				"sender": {
					"href": "/senders/some-sender-id"
				}
			}
		}`))
	})

	Context("when the list is empty", func() {
		It("returns an empty list (not null)", func() {
			response = campaigntypes.NewCampaignTypesListResponse("some-sender-id", []collections.CampaignType{})

			output, err := json.Marshal(response)
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(MatchJSON(`{
				"campaign_types": [],
				"_links": {
					"self": {
						"href": "/senders/some-sender-id/campaign_types"
					},
					"sender": {
						"href": "/senders/some-sender-id"
					}
				}
			}`))
		})
	})
})
