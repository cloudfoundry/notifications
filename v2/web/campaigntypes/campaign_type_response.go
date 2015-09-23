package campaigntypes

import (
	"fmt"

	"github.com/cloudfoundry-incubator/notifications/v2/collections"
)

type Link struct {
	Href string `json:"href"`
}

type CampaignTypeResponseLinks struct {
	Self Link `json:"self"`
}

type CampaignTypeResponse struct {
	ID          string                    `json:"id"`
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	Critical    bool                      `json:"critical"`
	TemplateID  string                    `json:"template_id"`
	Links       CampaignTypeResponseLinks `json:"_links"`
}

func NewCampaignTypeResponse(campaignType collections.CampaignType) CampaignTypeResponse {
	return CampaignTypeResponse{
		ID:          campaignType.ID,
		Name:        campaignType.Name,
		Description: campaignType.Description,
		Critical:    campaignType.Critical,
		TemplateID:  campaignType.TemplateID,
		Links: CampaignTypeResponseLinks{
			Self: Link{Href: fmt.Sprintf("/campaign_types/%s", campaignType.ID)},
		},
	}
}
