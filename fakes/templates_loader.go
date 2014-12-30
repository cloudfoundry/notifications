package fakes

import "github.com/cloudfoundry-incubator/notifications/postal"

type TemplatesLoader struct {
	Templates postal.Templates
	LoadError error
}

func NewTemplatesLoader() *TemplatesLoader {
	return &TemplatesLoader{}
}

func (fake *TemplatesLoader) LoadTemplates(clientID, kindID string) (postal.Templates, error) {
	return fake.Templates, fake.LoadError
}
