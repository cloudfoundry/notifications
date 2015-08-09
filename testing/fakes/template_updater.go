package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type TemplateUpdater struct {
	UpdateCall struct {
		Receives struct {
			Database   models.DatabaseInterface
			TemplateID string
			Template   models.Template
		}
		Returns struct {
			Error error
		}
	}
}

func NewTemplateUpdater() *TemplateUpdater {
	return &TemplateUpdater{}
}

func (tu *TemplateUpdater) Update(database models.DatabaseInterface, templateID string, template models.Template) error {
	tu.UpdateCall.Receives.Database = database
	tu.UpdateCall.Receives.TemplateID = templateID
	tu.UpdateCall.Receives.Template = template

	return tu.UpdateCall.Returns.Error
}
