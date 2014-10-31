package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type FakeTemplateFinder struct {
    TemplateName string
    FindError    error
    Templates    map[string]models.Template
}

func NewFakeTemplateFinder() *FakeTemplateFinder {
    return &FakeTemplateFinder{
        Templates: map[string]models.Template{},
    }
}

func (fake *FakeTemplateFinder) Find(templateName string) (models.Template, error) {
    fake.TemplateName = templateName
    return fake.Templates[templateName], fake.FindError
}
