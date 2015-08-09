package fakes

import "github.com/cloudfoundry-incubator/notifications/v2/collections"

type CampaignTypesCollection struct {
	SetCall struct {
		Receives struct {
			CampaignType collections.CampaignType
			Conn         collections.ConnectionInterface
		}
		Returns struct {
			CampaignType collections.CampaignType
			Err          error
		}
		WasCalled bool
	}

	ListCall struct {
		Receives struct {
			Conn     collections.ConnectionInterface
			SenderID string
			ClientID string
		}
		Returns struct {
			CampaignTypeList []collections.CampaignType
			Err              error
		}
	}

	GetCall struct {
		Receives struct {
			Conn           collections.ConnectionInterface
			CampaignTypeID string
			SenderID       string
			ClientID       string
		}
		Returns struct {
			CampaignType collections.CampaignType
			Err          error
		}
	}

	DeleteCall struct {
		Receives struct {
			Conn           collections.ConnectionInterface
			CampaignTypeID string
			SenderID       string
			ClientID       string
		}
		Returns struct {
			Err error
		}
	}
}

func NewCampaignTypesCollection() *CampaignTypesCollection {
	return &CampaignTypesCollection{}
}

func (c *CampaignTypesCollection) Set(conn collections.ConnectionInterface, campaignType collections.CampaignType, clientID string) (collections.CampaignType, error) {
	c.SetCall.WasCalled = true
	c.SetCall.Receives.Conn = conn
	c.SetCall.Receives.CampaignType = campaignType

	return c.SetCall.Returns.CampaignType, c.SetCall.Returns.Err
}

func (c *CampaignTypesCollection) List(conn collections.ConnectionInterface, senderID, clientID string) ([]collections.CampaignType, error) {
	c.ListCall.Receives.Conn = conn
	c.ListCall.Receives.SenderID = senderID
	c.ListCall.Receives.ClientID = clientID

	return c.ListCall.Returns.CampaignTypeList, c.ListCall.Returns.Err
}

func (c *CampaignTypesCollection) Get(conn collections.ConnectionInterface, campaignTypeID, senderID, clientID string) (collections.CampaignType, error) {
	c.GetCall.Receives.Conn = conn
	c.GetCall.Receives.CampaignTypeID = campaignTypeID
	c.GetCall.Receives.SenderID = senderID
	c.GetCall.Receives.ClientID = clientID

	return c.GetCall.Returns.CampaignType, c.GetCall.Returns.Err
}

func (c *CampaignTypesCollection) Delete(conn collections.ConnectionInterface, campaignTypeID, senderID, clientID string) error {
	c.DeleteCall.Receives.CampaignTypeID = campaignTypeID
	c.DeleteCall.Receives.SenderID = senderID
	c.DeleteCall.Receives.ClientID = clientID
	c.DeleteCall.Receives.Conn = conn

	return c.DeleteCall.Returns.Err
}
