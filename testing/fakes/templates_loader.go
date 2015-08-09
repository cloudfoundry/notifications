package fakes

import "github.com/cloudfoundry-incubator/notifications/postal"

type TemplatesLoader struct {
	LoadTemplatesCall struct {
		Receives struct {
			ClientID string
			KindID   string
		}
		Returns struct {
			Templates postal.Templates
			Error     error
		}
	}
}

func NewTemplatesLoader() *TemplatesLoader {
	return &TemplatesLoader{}
}

func (tl *TemplatesLoader) LoadTemplates(clientID, kindID string) (postal.Templates, error) {
	tl.LoadTemplatesCall.Receives.ClientID = clientID
	tl.LoadTemplatesCall.Receives.KindID = kindID

	return tl.LoadTemplatesCall.Returns.Templates, tl.LoadTemplatesCall.Returns.Error
}
