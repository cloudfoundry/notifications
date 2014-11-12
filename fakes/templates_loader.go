package fakes

import "github.com/cloudfoundry-incubator/notifications/postal"

type TemplatesLoader struct {
	ContentSuffix string
	Templates     postal.Templates
	LoadError     error
}

func NewTemplatesLoader() *TemplatesLoader {
	return &TemplatesLoader{}
}

func (fake *TemplatesLoader) LoadTemplates(subjectSuffix, contentSuffix, clientID, kindID string) (postal.Templates, error) {
	fake.ContentSuffix = contentSuffix
	return fake.Templates, fake.LoadError
}
