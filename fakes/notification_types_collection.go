package fakes

import (
	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/models"
)

type CampaignTypesCollection struct {
	AddCall struct {
		NotificationType       collections.CampaignType
		Conn                   models.ConnectionInterface
		ReturnNotificationType collections.CampaignType
		Err                    error
	}
	ListCall struct {
		ReturnCampaignTypeList []collections.CampaignType
		Err                        error
	}
	GetCall struct {
		ReturnCampaignType collections.CampaignType
		Err                    error
	}
}

func NewCampaignTypesCollection() *CampaignTypesCollection {
	return &CampaignTypesCollection{}
}

func (c *CampaignTypesCollection) Add(conn models.ConnectionInterface, notificationType collections.CampaignType, clientID string) (collections.CampaignType, error) {
	c.AddCall.Conn = conn
	c.AddCall.NotificationType = notificationType
	return c.AddCall.ReturnNotificationType, c.AddCall.Err
}

func (c *CampaignTypesCollection) List(conn models.ConnectionInterface, senderID, clientID string) ([]collections.CampaignType, error) {
	return c.ListCall.ReturnCampaignTypeList, c.ListCall.Err
}

func (c *CampaignTypesCollection) Get(conn models.ConnectionInterface, notificationTypeID, senderID, clientID string) (collections.CampaignType, error) {
	return c.GetCall.ReturnCampaignType, c.GetCall.Err
}
