package mocks

import "github.com/cloudfoundry-incubator/notifications/v1/collections"

type TemplateDeleter struct {
	DeleteCall struct {
		Receives struct {
			Connection collections.ConnectionInterface
			TemplateID string
		}
		Returns struct {
			Error error
		}
	}
}

func NewTemplateDeleter() *TemplateDeleter {
	return &TemplateDeleter{}
}

func (td *TemplateDeleter) Delete(connection collections.ConnectionInterface, templateID string) error {
	td.DeleteCall.Receives.Connection = connection
	td.DeleteCall.Receives.TemplateID = templateID

	return td.DeleteCall.Returns.Error
}
