package mocks

import "github.com/cloudfoundry-incubator/notifications/v2/collections"

type CampaignStatusesCollection struct {
	GetCall struct {
		Receives struct {
			Connection collections.ConnectionInterface
			CampaignID string
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

func (csc *CampaignStatusesCollection) Get(conn collections.ConnectionInterface, campaignID string) (collections.CampaignStatus, error) {
	csc.GetCall.Receives.Connection = conn
	csc.GetCall.Receives.CampaignID = campaignID

	return csc.GetCall.Returns.CampaignStatus, csc.GetCall.Returns.Error
}
