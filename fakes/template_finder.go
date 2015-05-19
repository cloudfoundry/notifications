package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type TemplateFinder struct {
	Templates map[string]models.Template

	FindByIDCall struct {
		Arguments []interface{}
		Error     error
	}
}

func NewTemplateFinder() *TemplateFinder {
	return &TemplateFinder{
		Templates: map[string]models.Template{},
	}
}

func (fake *TemplateFinder) FindByID(database models.DatabaseInterface, templateID string) (models.Template, error) {
	fake.FindByIDCall.Arguments = []interface{}{database, templateID}

	return fake.Templates[templateID], fake.FindByIDCall.Error
}
