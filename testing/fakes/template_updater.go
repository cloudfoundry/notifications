package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type TemplateUpdater struct {
	UpdateCall struct {
		Arguments []interface{}
		Error     error
	}
}

func NewTemplateUpdater() *TemplateUpdater {
	return &TemplateUpdater{}
}

func (fake *TemplateUpdater) Update(database models.DatabaseInterface, templateID string, template models.Template) error {
	fake.UpdateCall.Arguments = []interface{}{database, templateID, template}

	return fake.UpdateCall.Error
}
