package fakes

import "github.com/cloudfoundry-incubator/notifications/v2/collections"

type SendersCollection struct {
	SetCall struct {
		Receives struct {
			Conn   collections.ConnectionInterface
			Sender collections.Sender
		}
		Returns struct {
			Sender collections.Sender
			Err    error
		}
	}

	ListCall struct {
		Receives struct {
			Conn     collections.ConnectionInterface
			ClientID string
		}
		Returns struct {
			SenderList []collections.Sender
			Err        error
		}
	}

	GetCall struct {
		Receives struct {
			Conn     collections.ConnectionInterface
			SenderID string
			ClientID string
		}
		Returns struct {
			Sender collections.Sender
			Err    error
		}
	}
}

func NewSendersCollection() *SendersCollection {
	return &SendersCollection{}
}

func (c *SendersCollection) Set(conn collections.ConnectionInterface, sender collections.Sender) (collections.Sender, error) {
	c.SetCall.Receives.Conn = conn
	c.SetCall.Receives.Sender = sender

	return c.SetCall.Returns.Sender, c.SetCall.Returns.Err
}

func (c *SendersCollection) Get(conn collections.ConnectionInterface, senderID, clientID string) (collections.Sender, error) {
	c.GetCall.Receives.Conn = conn
	c.GetCall.Receives.SenderID = senderID
	c.GetCall.Receives.ClientID = clientID

	return c.GetCall.Returns.Sender, c.GetCall.Returns.Err
}

func (c *SendersCollection) List(conn collections.ConnectionInterface, clientID string) ([]collections.Sender, error) {
	c.ListCall.Receives.Conn = conn
	c.ListCall.Receives.ClientID = clientID

	return c.ListCall.Returns.SenderList, c.ListCall.Returns.Err
}
