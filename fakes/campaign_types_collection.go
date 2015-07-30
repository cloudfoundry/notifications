package fakes

import (
	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/models"
)

type CampaignTypesCollection struct {
	SetCall struct {
		CampaignType       collections.CampaignType
		Conn               models.ConnectionInterface
		ReturnCampaignType collections.CampaignType
		Err                error
	}

	ListCall struct {
		ReturnCampaignTypeList []collections.CampaignType
		Err                    error
	}

	GetCall struct {
		ReturnCampaignType collections.CampaignType
		Err                error
	}
}

func NewCampaignTypesCollection() *CampaignTypesCollection {
	return &CampaignTypesCollection{}
}

func (c *CampaignTypesCollection) Set(conn models.ConnectionInterface, campaignType collections.CampaignType, clientID string) (collections.CampaignType, error) {
	c.SetCall.Conn = conn
	c.SetCall.CampaignType = campaignType
	return c.SetCall.ReturnCampaignType, c.SetCall.Err
}

func (c *CampaignTypesCollection) List(conn models.ConnectionInterface, senderID, clientID string) ([]collections.CampaignType, error) {
	return c.ListCall.ReturnCampaignTypeList, c.ListCall.Err
}

func (c *CampaignTypesCollection) Get(conn models.ConnectionInterface, campaignTypeID, senderID, clientID string) (collections.CampaignType, error) {
	return c.GetCall.ReturnCampaignType, c.GetCall.Err
}
