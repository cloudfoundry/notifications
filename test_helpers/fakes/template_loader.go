package fakes

import "github.com/cloudfoundry-incubator/notifications/postal"

type FakeTemplateLoader struct {
    Templates postal.Templates
}

func (fake *FakeTemplateLoader) LoadNamedTemplates(subjectTemplateName, textTemplateName, htmlTemplateName string) (postal.Templates, error) {
    return fake.Templates, nil
}

func (fake *FakeTemplateLoader) LoadNamedTemplatesWithClientAndKind(subjectTemplateName, textTemplateName, htmlTemplateName, clientID, kind string) (postal.Templates, error) {
    return fake.Templates, nil
}
