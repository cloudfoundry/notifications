package mocks

import "github.com/cloudfoundry-incubator/notifications/db"

type TemplateDeleter struct {
	DeleteCall struct {
		Receives struct {
			Database   db.DatabaseInterface
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

func (td *TemplateDeleter) Delete(database db.DatabaseInterface, templateID string) error {
	td.DeleteCall.Receives.Database = database
	td.DeleteCall.Receives.TemplateID = templateID

	return td.DeleteCall.Returns.Error
}
