package fakes

import (
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
)

type TemplatesCollection struct {
	SetCall struct {
		Receives struct {
			Conn     db.ConnectionInterface
			Template collections.Template
			ClientID string
		}
		Returns struct {
			Template collections.Template
			Err      error
		}
	}

	GetCall struct {
		Receives struct {
			Conn       collections.ConnectionInterface
			TemplateID string
			ClientID   string
		}
		Returns struct {
			Template collections.Template
			Err      error
		}
	}
}

func NewTemplatesCollection() *TemplatesCollection {
	return &TemplatesCollection{}
}

func (c *TemplatesCollection) Set(conn db.ConnectionInterface, template collections.Template, clientID string) (collections.Template, error) {
	c.SetCall.Receives.Conn = conn
	c.SetCall.Receives.Template = template
	c.SetCall.Receives.ClientID = clientID

	return c.SetCall.Returns.Template, c.SetCall.Returns.Err
}

func (c *TemplatesCollection) Get(conn collections.ConnectionInterface, templateID, clientID string) (collections.Template, error) {
	c.GetCall.Receives.Conn = conn
	c.GetCall.Receives.TemplateID = templateID
	c.GetCall.Receives.ClientID = clientID

	return c.GetCall.Returns.Template, c.GetCall.Returns.Err
}
