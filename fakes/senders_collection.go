package fakes

import (
	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/models"
)

type SendersCollection struct {
	AddCall struct {
		Conn         models.ConnectionInterface
		Sender       collections.Sender
		ReturnSender collections.Sender
		Err          error
	}
}

func NewSendersCollection() *SendersCollection {
	return &SendersCollection{}
}

func (c *SendersCollection) Add(conn models.ConnectionInterface, sender collections.Sender) (collections.Sender, error) {
	c.AddCall.Conn = conn
	c.AddCall.Sender = sender
	return c.AddCall.ReturnSender, c.AddCall.Err
}
