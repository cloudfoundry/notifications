package campaigntypes

import (
	"fmt"

	"github.com/cloudfoundry-incubator/notifications/v2/collections"
)

type CampaignTypesListResponse struct {
	CampaignTypes []CampaignTypeResponse         `json:"campaign_types"`
	Links         CampaignTypesListResponseLinks `json:"_links"`
}

type CampaignTypesListResponseLinks struct {
	Self   Link `json:"self"`
	Sender Link `json:"sender"`
}

func NewCampaignTypesListResponse(senderID string, campaignTypes []collections.CampaignType) CampaignTypesListResponse {
	campaignTypeResponseList := []CampaignTypeResponse{}

	for _, campaignType := range campaignTypes {
		campaignTypeResponseList = append(campaignTypeResponseList, NewCampaignTypeResponse(campaignType))
	}

	return CampaignTypesListResponse{
		CampaignTypes: campaignTypeResponseList,
		Links: CampaignTypesListResponseLinks{
			Self:   Link{fmt.Sprintf("/senders/%s/campaign_types", senderID)},
			Sender: Link{fmt.Sprintf("/senders/%s", senderID)},
		},
	}
}
