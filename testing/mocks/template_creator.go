package mocks

import (
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
)

type TemplateCreator struct {
	CreateCall struct {
		Receives struct {
			Database services.DatabaseInterface
			Template models.Template
		}
		Returns struct {
			TemplateGUID string
			Error        error
		}
	}
}

func NewTemplateCreator() *TemplateCreator {
	return &TemplateCreator{}
}

func (tc *TemplateCreator) Create(database services.DatabaseInterface, template models.Template) (string, error) {
	tc.CreateCall.Receives.Database = database
	tc.CreateCall.Receives.Template = template

	return tc.CreateCall.Returns.TemplateGUID, tc.CreateCall.Returns.Error
}
