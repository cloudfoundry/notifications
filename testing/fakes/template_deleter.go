package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type TemplateDeleter struct {
	DeleteCall struct {
		Arguments []interface{}
		Error     error
	}
}

func NewTemplateDeleter() *TemplateDeleter {
	return &TemplateDeleter{}
}

func (fake *TemplateDeleter) Delete(database models.DatabaseInterface, templateID string) error {
	fake.DeleteCall.Arguments = []interface{}{database, templateID}
	return fake.DeleteCall.Error
}
