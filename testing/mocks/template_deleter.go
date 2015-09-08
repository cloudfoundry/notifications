package mocks

import "github.com/cloudfoundry-incubator/notifications/v1/services"

type TemplateDeleter struct {
	DeleteCall struct {
		Receives struct {
			Database   services.DatabaseInterface
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

func (td *TemplateDeleter) Delete(database services.DatabaseInterface, templateID string) error {
	td.DeleteCall.Receives.Database = database
	td.DeleteCall.Receives.TemplateID = templateID

	return td.DeleteCall.Returns.Error
}
