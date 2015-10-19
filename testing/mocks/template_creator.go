package mocks

import "github.com/cloudfoundry-incubator/notifications/v1/collections"

type TemplateCreator struct {
	CreateCall struct {
		Receives struct {
			Connection collections.ConnectionInterface
			Template   collections.Template
		}
		Returns struct {
			Template collections.Template
			Error    error
		}
	}
}

func NewTemplateCreator() *TemplateCreator {
	return &TemplateCreator{}
}

func (tc *TemplateCreator) Create(connection collections.ConnectionInterface, template collections.Template) (collections.Template, error) {
	tc.CreateCall.Receives.Connection = connection
	tc.CreateCall.Receives.Template = template

	return tc.CreateCall.Returns.Template, tc.CreateCall.Returns.Error
}
