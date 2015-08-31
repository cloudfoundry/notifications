package mocks

import "github.com/cloudfoundry-incubator/notifications/v2/collections"

type SendersCollection struct {
	SetCall struct {
		Receives struct {
			Connection collections.ConnectionInterface
			Sender     collections.Sender
		}
		Returns struct {
			Sender collections.Sender
			Error  error
		}
	}

	ListCall struct {
		Receives struct {
			Connection collections.ConnectionInterface
			ClientID   string
		}
		Returns struct {
			SenderList []collections.Sender
			Error      error
		}
	}

	GetCall struct {
		Receives struct {
			Connection collections.ConnectionInterface
			SenderID   string
			ClientID   string
		}
		Returns struct {
			Sender collections.Sender
			Error  error
		}
	}

	DeleteCall struct {
		Receives struct {
			Connection collections.ConnectionInterface
			SenderID   string
			ClientID   string
		}
		Returns struct {
			Error error
		}
	}
}

func NewSendersCollection() *SendersCollection {
	return &SendersCollection{}
}

func (c *SendersCollection) Set(conn collections.ConnectionInterface, sender collections.Sender) (collections.Sender, error) {
	c.SetCall.Receives.Connection = conn
	c.SetCall.Receives.Sender = sender

	return c.SetCall.Returns.Sender, c.SetCall.Returns.Error
}

func (c *SendersCollection) Get(conn collections.ConnectionInterface, senderID, clientID string) (collections.Sender, error) {
	c.GetCall.Receives.Connection = conn
	c.GetCall.Receives.SenderID = senderID
	c.GetCall.Receives.ClientID = clientID

	return c.GetCall.Returns.Sender, c.GetCall.Returns.Error
}

func (c *SendersCollection) List(conn collections.ConnectionInterface, clientID string) ([]collections.Sender, error) {
	c.ListCall.Receives.Connection = conn
	c.ListCall.Receives.ClientID = clientID

	return c.ListCall.Returns.SenderList, c.ListCall.Returns.Error
}

func (c *SendersCollection) Delete(conn collections.ConnectionInterface, senderID, clientID string) error {
	c.DeleteCall.Receives.Connection = conn
	c.DeleteCall.Receives.SenderID = senderID
	c.DeleteCall.Receives.ClientID = clientID
	return c.DeleteCall.Returns.Error
}
