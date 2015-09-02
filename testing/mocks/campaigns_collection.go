package mocks

import "github.com/cloudfoundry-incubator/notifications/v2/collections"

type CampaignsCollection struct {
	CreateCall struct {
		Receives struct {
			Connection       collections.ConnectionInterface
			Campaign         collections.Campaign
			ClientID         string
			HasCriticalScope bool
		}
		Returns struct {
			Campaign collections.Campaign
			Error    error
		}
		WasCalled bool
	}

	GetCall struct {
		Receives struct {
			Connection collections.ConnectionInterface
			CampaignID string
			SenderID   string
			ClientID   string
		}
		Returns struct {
			Campaign collections.Campaign
			Error    error
		}
	}
}

func NewCampaignsCollection() *CampaignsCollection {
	return &CampaignsCollection{}
}

func (c *CampaignsCollection) Create(conn collections.ConnectionInterface, campaign collections.Campaign, clientID string, hasCriticalScope bool) (collections.Campaign, error) {
	c.CreateCall.Receives.Connection = conn
	c.CreateCall.Receives.Campaign = campaign
	c.CreateCall.Receives.ClientID = clientID
	c.CreateCall.Receives.HasCriticalScope = hasCriticalScope
	c.CreateCall.WasCalled = true

	return c.CreateCall.Returns.Campaign, c.CreateCall.Returns.Error
}

func (c *CampaignsCollection) Get(connection collections.ConnectionInterface, campaignID, senderID, clientID string) (collections.Campaign, error) {
	c.GetCall.Receives.Connection = connection
	c.GetCall.Receives.CampaignID = campaignID
	c.GetCall.Receives.SenderID = senderID
	c.GetCall.Receives.ClientID = clientID

	return c.GetCall.Returns.Campaign, c.GetCall.Returns.Error
}
