package mocks

import "github.com/cloudfoundry-incubator/notifications/v2/collections"

type CampaignsCollection struct {
	CreateCall struct {
		Receives struct {
			Conn             collections.ConnectionInterface
			Campaign         collections.Campaign
			HasCriticalScope bool
		}
		Returns struct {
			Campaign collections.Campaign
			Error    error
		}
		WasCalled bool
	}
}

func NewCampaignsCollection() *CampaignsCollection {
	return &CampaignsCollection{}
}

func (c *CampaignsCollection) Create(conn collections.ConnectionInterface, campaign collections.Campaign, hasCriticalScope bool) (collections.Campaign, error) {
	c.CreateCall.Receives.Conn = conn
	c.CreateCall.Receives.Campaign = campaign
	c.CreateCall.Receives.HasCriticalScope = hasCriticalScope
	c.CreateCall.WasCalled = true

	return c.CreateCall.Returns.Campaign, c.CreateCall.Returns.Error
}
