package mocks

import "github.com/cloudfoundry-incubator/notifications/v2/collections"

type CampaignStatusesCollection struct {
	GetCall struct {
		Receives struct {
			Connection collections.ConnectionInterface
			CampaignID string
			ClientID   string
		}
		Returns struct {
			CampaignStatus collections.CampaignStatus
			Error          error
		}
	}
}

func NewCampaignStatusesCollection() *CampaignStatusesCollection {
	return &CampaignStatusesCollection{}
}

func (csc *CampaignStatusesCollection) Get(conn collections.ConnectionInterface, campaignID, clientID string) (collections.CampaignStatus, error) {
	csc.GetCall.Receives.Connection = conn
	csc.GetCall.Receives.CampaignID = campaignID
	csc.GetCall.Receives.ClientID = clientID

	return csc.GetCall.Returns.CampaignStatus, csc.GetCall.Returns.Error
}
