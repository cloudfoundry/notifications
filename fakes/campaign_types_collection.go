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
		WasCalled          bool
	}

	ListCall struct {
		ReturnCampaignTypeList []collections.CampaignType
		Conn                   models.ConnectionInterface
		SenderID               string
		ClientID               string
		Err                    error
	}

	GetCall struct {
		ReturnCampaignType collections.CampaignType
		Conn               models.ConnectionInterface
		CampaignTypeID     string
		SenderID           string
		ClientID           string
		Err                error
	}

	DeleteCall struct {
		CampaignTypeID string
		SenderID       string
		ClientID       string
		Conn           models.ConnectionInterface
		Err            error
	}
}

func NewCampaignTypesCollection() *CampaignTypesCollection {
	return &CampaignTypesCollection{}
}

func (c *CampaignTypesCollection) Set(conn models.ConnectionInterface, campaignType collections.CampaignType, clientID string) (collections.CampaignType, error) {
	c.SetCall.WasCalled = true
	c.SetCall.Conn = conn
	c.SetCall.CampaignType = campaignType
	return c.SetCall.ReturnCampaignType, c.SetCall.Err
}

func (c *CampaignTypesCollection) List(conn models.ConnectionInterface, senderID, clientID string) ([]collections.CampaignType, error) {
	c.ListCall.Conn = conn
	c.ListCall.SenderID = senderID
	c.ListCall.ClientID = clientID
	return c.ListCall.ReturnCampaignTypeList, c.ListCall.Err
}

func (c *CampaignTypesCollection) Get(conn models.ConnectionInterface, campaignTypeID, senderID, clientID string) (collections.CampaignType, error) {
	c.GetCall.Conn = conn
	c.GetCall.CampaignTypeID = campaignTypeID
	c.GetCall.SenderID = senderID
	c.GetCall.ClientID = clientID
	return c.GetCall.ReturnCampaignType, c.GetCall.Err
}

func (c *CampaignTypesCollection) Delete(conn models.ConnectionInterface, campaignTypeID, senderID, clientID string) error {
	c.DeleteCall.CampaignTypeID = campaignTypeID
	c.DeleteCall.SenderID = senderID
	c.DeleteCall.ClientID = clientID
	c.DeleteCall.Conn = conn

	return c.DeleteCall.Err
}
