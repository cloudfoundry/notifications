package mocks

import "github.com/cloudfoundry-incubator/notifications/postal/common"

type TemplatesLoader struct {
	LoadTemplatesCall struct {
		Receives struct {
			ClientID   string
			KindID     string
			TemplateID string
		}
		Returns struct {
			Templates common.Templates
			Error     error
		}
	}
}

func NewTemplatesLoader() *TemplatesLoader {
	return &TemplatesLoader{}
}

func (tl *TemplatesLoader) LoadTemplates(clientID, kindID, templateID string) (common.Templates, error) {
	tl.LoadTemplatesCall.Receives.ClientID = clientID
	tl.LoadTemplatesCall.Receives.KindID = kindID
	tl.LoadTemplatesCall.Receives.TemplateID = templateID

	return tl.LoadTemplatesCall.Returns.Templates, tl.LoadTemplatesCall.Returns.Error
}
