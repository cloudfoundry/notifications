package collections

import "github.com/cloudfoundry-incubator/notifications/models"

type Template struct {
	ID       string
	Name     string
	Html     string
	Text     string
	Subject  string
	Metadata string
}

type TemplatesCollection struct {
}

func NewTemplatesCollection() TemplatesCollection {
	return TemplatesCollection{}
}

func (c TemplatesCollection) Get(conn models.ConnectionInterface, templateID, clientID string) (Template, error) {
	return Template{}, nil
}
