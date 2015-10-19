package campaigns

import (
	"fmt"

	"github.com/cloudfoundry-incubator/notifications/v2/collections"
)

type Link struct {
	Href string `json:"href"`
}

type CampaignResponseLinks struct {
	Self         Link `json:"self"`
	Template     Link `json:"template"`
	CampaignType Link `json:"campaign_type"`
	Status       Link `json:"status"`
}

type CampaignResponse struct {
	ID             string                `json:"id"`
	SendTo         map[string][]string   `json:"send_to"`
	CampaignTypeID string                `json:"campaign_type_id"`
	Text           string                `json:"text"`
	HTML           string                `json:"html"`
	Subject        string                `json:"subject"`
	TemplateID     string                `json:"template_id"`
	ReplyTo        string                `json:"reply_to"`
	Links          CampaignResponseLinks `json:"_links"`
}

func NewCampaignResponse(campaign collections.Campaign) CampaignResponse {
	return CampaignResponse{
		ID:             campaign.ID,
		SendTo:         campaign.SendTo,
		CampaignTypeID: campaign.CampaignTypeID,
		Text:           campaign.Text,
		HTML:           campaign.HTML,
		Subject:        campaign.Subject,
		TemplateID:     campaign.TemplateID,
		ReplyTo:        campaign.ReplyTo,
		Links: CampaignResponseLinks{
			Self:         Link{fmt.Sprintf("/campaigns/%s", campaign.ID)},
			Template:     Link{fmt.Sprintf("/templates/%s", campaign.TemplateID)},
			CampaignType: Link{fmt.Sprintf("/campaign_types/%s", campaign.CampaignTypeID)},
			Status:       Link{fmt.Sprintf("/campaigns/%s/status", campaign.ID)},
		},
	}
}
