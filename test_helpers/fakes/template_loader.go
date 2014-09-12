package fakes

import "github.com/cloudfoundry-incubator/notifications/postal"

type FakeTemplateLoader struct {
    Templates postal.Templates
}

func (fake *FakeTemplateLoader) Load(subject string, guid postal.TypedGUID, clientID string, kind string) (postal.Templates, error) {
    return fake.Templates, nil
}
