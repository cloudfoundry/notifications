package fakes

import (
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
)

type TemplateUpdater struct {
	UpdateCall struct {
		Receives struct {
			Database   db.DatabaseInterface
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

func (tu *TemplateUpdater) Update(database db.DatabaseInterface, templateID string, template models.Template) error {
	tu.UpdateCall.Receives.Database = database
	tu.UpdateCall.Receives.TemplateID = templateID
	tu.UpdateCall.Receives.Template = template

	return tu.UpdateCall.Returns.Error
}
