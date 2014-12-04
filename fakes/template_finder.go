package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type TemplateFinder struct {
	TemplateID string
	FindError  error
	Templates  map[string]models.Template
}

func NewTemplateFinder() *TemplateFinder {
	return &TemplateFinder{
		Templates: map[string]models.Template{},
	}
}

func (fake *TemplateFinder) Find(templateID string) (models.Template, error) {
	fake.TemplateID = templateID
	return fake.Templates[templateID], fake.FindError
}
