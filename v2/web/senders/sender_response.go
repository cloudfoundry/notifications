package senders

import (
	"fmt"

	"github.com/cloudfoundry-incubator/notifications/v2/collections"
)

type SenderResponse struct {
	ID    string              `json:"id"`
	Name  string              `json:"name"`
	Links SenderResponseLinks `json:"_links"`
}

type SenderResponseLinks struct {
	Self          Link `json:"self"`
	CampaignTypes Link `json:"campaign_types"`
	Campaigns     Link `json:"campaigns"`
}

type Link struct {
	Href string `json:"href"`
}

func NewSenderResponse(sender collections.Sender) SenderResponse {
	return SenderResponse{
		ID:   sender.ID,
		Name: sender.Name,
		Links: SenderResponseLinks{
			Self:          Link{fmt.Sprintf("/senders/%s", sender.ID)},
			CampaignTypes: Link{fmt.Sprintf("/senders/%s/campaign_types", sender.ID)},
			Campaigns:     Link{fmt.Sprintf("/senders/%s/campaigns", sender.ID)},
		},
	}
}
