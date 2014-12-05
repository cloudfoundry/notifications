package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type TemplateCreator struct {
	CreateArgument models.Template
	CreateError    error
}

func NewTemplateCreator() *TemplateCreator {
	return &TemplateCreator{}
}

func (fake *TemplateCreator) Create(template models.Template) (string, error) {
	fake.CreateArgument = template
	return "guid", fake.CreateError
}
