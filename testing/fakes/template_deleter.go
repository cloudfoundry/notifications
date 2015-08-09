package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type TemplateDeleter struct {
	DeleteCall struct {
		Receives struct {
			Database   models.DatabaseInterface
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

func (td *TemplateDeleter) Delete(database models.DatabaseInterface, templateID string) error {
	td.DeleteCall.Receives.Database = database
	td.DeleteCall.Receives.TemplateID = templateID

	return td.DeleteCall.Returns.Error
}
