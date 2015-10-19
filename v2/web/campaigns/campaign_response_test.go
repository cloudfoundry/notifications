package campaigns_test

import (
	"encoding/json"

	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/web/campaigns"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CampaignResponse", func() {
	It("provides a JSON representation of a campaign resource", func() {
		campaign := collections.Campaign{
			ID: "some-campaign-id",
			SendTo: map[string][]string{
				"emails": {"me@example.com"},
				"users":  {"some-user-guid"},
				"spaces": {"some-space-guid"},
				"orgs":   {"some-org-guid"},
			},
			CampaignTypeID: "some-campaign-type-id",
			Text:           "some-text",
			HTML:           "some-html",
			Subject:        "some-subject",
			TemplateID:     "some-template-id",
			ReplyTo:        "some-reply-to",
		}

		response := campaigns.NewCampaignResponse(campaign)
		Expect(response).To(Equal(campaigns.CampaignResponse{
			ID: "some-campaign-id",
			SendTo: map[string][]string{
				"emails": {"me@example.com"},
				"users":  {"some-user-guid"},
				"spaces": {"some-space-guid"},
				"orgs":   {"some-org-guid"},
			},
			CampaignTypeID: "some-campaign-type-id",
			Text:           "some-text",
			HTML:           "some-html",
			Subject:        "some-subject",
			TemplateID:     "some-template-id",
			ReplyTo:        "some-reply-to",
			Links: campaigns.CampaignResponseLinks{
				Self:         campaigns.Link{"/campaigns/some-campaign-id"},
				Template:     campaigns.Link{"/templates/some-template-id"},
				CampaignType: campaigns.Link{"/campaign_types/some-campaign-type-id"},
				Status:       campaigns.Link{"/campaigns/some-campaign-id/status"},
			},
		}))
	})

	It("can marshal into JSON", func() {
		campaign := collections.Campaign{
			ID: "some-campaign-id",
			SendTo: map[string][]string{
				"emails": {"me@example.com"},
				"users":  {"some-user-guid"},
				"spaces": {"some-space-guid"},
				"orgs":   {"some-org-guid"},
			},
			CampaignTypeID: "some-campaign-type-id",
			Text:           "some-text",
			HTML:           "some-html",
			Subject:        "some-subject",
			TemplateID:     "some-template-id",
			ReplyTo:        "some-reply-to",
		}

		output, err := json.Marshal(campaigns.NewCampaignResponse(campaign))
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(MatchJSON(`{
			"id":   "some-campaign-id",
			"send_to": {
				"emails": ["me@example.com"],
				"users":  ["some-user-guid"],
				"spaces": ["some-space-guid"],
				"orgs":   ["some-org-guid"]
			},
			"campaign_type_id": "some-campaign-type-id",
			"text": "some-text",
			"html": "some-html",
			"subject": "some-subject",
			"template_id": "some-template-id",
			"reply_to": "some-reply-to",
			"_links": {
				"self": {
					"href": "/campaigns/some-campaign-id"
				},
				"template": {
					"href": "/templates/some-template-id"
				},
				"campaign_type": {
					"href": "/campaign_types/some-campaign-type-id"
				},
				"status": {
					"href": "/campaigns/some-campaign-id/status"
				}
			}
		}`))
	})
})
