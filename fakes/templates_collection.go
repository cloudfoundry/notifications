package fakes

import (
	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/models"
)

type TemplatesCollection struct {
	SetCall struct {
		Template       collections.Template
		ClientID       string
		Conn           models.ConnectionInterface
		ReturnTemplate collections.Template
		Err            error
	}

	GetCall struct {
		TemplateID     string
		ClientID       string
		Conn           models.ConnectionInterface
		ReturnTemplate collections.Template
		Err            error
	}
}

func NewTemplatesCollection() *TemplatesCollection {
	return &TemplatesCollection{}
}

func (c *TemplatesCollection) Set(conn models.ConnectionInterface, template collections.Template, clientID string) (collections.Template, error) {
	c.SetCall.Template = template
	c.SetCall.ClientID = clientID
	c.SetCall.Conn = conn
	return c.SetCall.ReturnTemplate, c.SetCall.Err
}

func (c *TemplatesCollection) Get(conn models.ConnectionInterface, templateID, clientID string) (collections.Template, error) {
	c.GetCall.TemplateID = templateID
	c.GetCall.ClientID = clientID
	c.GetCall.Conn = conn
	return c.GetCall.ReturnTemplate, c.GetCall.Err
}
