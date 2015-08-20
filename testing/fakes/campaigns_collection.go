package fakes

import "github.com/cloudfoundry-incubator/notifications/v2/collections"

type CampaignsCollection struct {
	CreateCall struct {
		Receives struct {
			Conn     collections.ConnectionInterface
			Campaign collections.Campaign
		}
		Returns struct {
			Campaign collections.Campaign
			Err      error
		}
		WasCalled bool
	}
}

func NewCampaignsCollection() *CampaignsCollection {
	return &CampaignsCollection{}
}

func (c *CampaignsCollection) Create(conn collections.ConnectionInterface, campaign collections.Campaign) (collections.Campaign, error) {
	c.CreateCall.Receives.Conn = conn
	c.CreateCall.Receives.Campaign = campaign
	c.CreateCall.WasCalled = true

	return c.CreateCall.Returns.Campaign, c.CreateCall.Returns.Err
}
