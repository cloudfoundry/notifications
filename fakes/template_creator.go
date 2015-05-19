package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type TemplateCreator struct {
	CreateCall struct {
		Arguments []interface{}
		Error     error
	}
}

func NewTemplateCreator() *TemplateCreator {
	return &TemplateCreator{}
}

func (fake *TemplateCreator) Create(database models.DatabaseInterface, template models.Template) (string, error) {
	fake.CreateCall.Arguments = []interface{}{database, template}
	return "guid", fake.CreateCall.Error
}
