package fakes

import "github.com/cloudfoundry-incubator/notifications/postal"

type TemplatesLoader struct {
	SubjectSuffix string
	ContentSuffix string
	Templates     postal.Templates
	LoadError     error
}

func NewTemplatesLoader() *TemplatesLoader {
	return &TemplatesLoader{}
}

func (fake *TemplatesLoader) DeprecatedLoadTemplates(subjectSuffix, contentSuffix, clientID, kindID string) (postal.Templates, error) {
	fake.SubjectSuffix = subjectSuffix
	fake.ContentSuffix = contentSuffix
	return fake.Templates, fake.LoadError
}

func (fake *TemplatesLoader) LoadTemplates(clientID, kindID, contentSuffix, subjectSuffix string) (postal.Templates, error) {
	fake.SubjectSuffix = subjectSuffix
	fake.ContentSuffix = contentSuffix
	return fake.Templates, fake.LoadError
}
