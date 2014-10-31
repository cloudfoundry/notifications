package fakes

import "github.com/cloudfoundry-incubator/notifications/postal"

type FakeTemplatesLoader struct {
    Templates postal.Templates
    LoadError error
}

func (fake *FakeTemplatesLoader) LoadTemplates(string, string, string, string) (postal.Templates, error) {
    return fake.Templates, fake.LoadError
}
