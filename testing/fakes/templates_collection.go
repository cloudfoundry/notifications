package fakes

import (
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
)

type TemplatesCollection struct {
	SetCall struct {
		Receives struct {
			Connection db.ConnectionInterface
			Template   collections.Template
		}
		Returns struct {
			Template collections.Template
			Err      error
		}
	}

	GetCall struct {
		Receives struct {
			Connection collections.ConnectionInterface
			TemplateID string
			ClientID   string
		}
		Returns struct {
			Template collections.Template
			Err      error
		}
	}

	DeleteCall struct {
		Receives struct {
			Connection collections.ConnectionInterface
			TemplateID string
		}
		Returns struct {
			Err error
		}
	}
}

func NewTemplatesCollection() *TemplatesCollection {
	return &TemplatesCollection{}
}

func (c *TemplatesCollection) Set(conn collections.ConnectionInterface, template collections.Template) (collections.Template, error) {
	c.SetCall.Receives.Connection = conn
	c.SetCall.Receives.Template = template

	return c.SetCall.Returns.Template, c.SetCall.Returns.Err
}

func (c *TemplatesCollection) Get(conn collections.ConnectionInterface, templateID, clientID string) (collections.Template, error) {
	c.GetCall.Receives.Connection = conn
	c.GetCall.Receives.TemplateID = templateID
	c.GetCall.Receives.ClientID = clientID

	return c.GetCall.Returns.Template, c.GetCall.Returns.Err
}

func (c *TemplatesCollection) Delete(conn collections.ConnectionInterface, templateID string) error {
	c.DeleteCall.Receives.Connection = conn
	c.DeleteCall.Receives.TemplateID = templateID

	return c.DeleteCall.Returns.Err
}
