package fakes

import "github.com/cloudfoundry-incubator/notifications/postal"

type FakeTemplatesLoader struct {
    ContentSuffix string
    Templates     postal.Templates
    LoadError     error
}

func (fake *FakeTemplatesLoader) LoadTemplates(subjectSuffix, contentSuffix, clientID, kindID string) (postal.Templates, error) {
    fake.ContentSuffix = contentSuffix
    return fake.Templates, fake.LoadError
}
