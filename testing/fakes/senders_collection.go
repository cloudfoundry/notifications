package fakes

import (
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
)

type SendersCollection struct {
	SetCall struct {
		Conn         models.ConnectionInterface
		Sender       collections.Sender
		ReturnSender collections.Sender
		Err          error
	}

	GetCall struct {
		Conn         models.ConnectionInterface
		SenderID     string
		ClientID     string
		ReturnSender collections.Sender
		Err          error
	}
}

func NewSendersCollection() *SendersCollection {
	return &SendersCollection{}
}

func (c *SendersCollection) Set(conn models.ConnectionInterface, sender collections.Sender) (collections.Sender, error) {
	c.SetCall.Conn = conn
	c.SetCall.Sender = sender
	return c.SetCall.ReturnSender, c.SetCall.Err
}

func (c *SendersCollection) Get(conn models.ConnectionInterface, senderID, clientID string) (collections.Sender, error) {
	c.GetCall.Conn = conn
	c.GetCall.SenderID = senderID
	c.GetCall.ClientID = clientID
	return c.GetCall.ReturnSender, c.GetCall.Err
}
