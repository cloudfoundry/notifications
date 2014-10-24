package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type FakeTemplateFinder struct {
    TemplateName     string
    TemplateResponse models.Template
    FindError        error
}

func NewFakeTemplateFinder(templateResponse models.Template) *FakeTemplateFinder {
    return &FakeTemplateFinder{
        TemplateResponse: templateResponse,
    }
}

func (fake *FakeTemplateFinder) Find(templateName string) (models.Template, error) {
    fake.TemplateName = templateName
    return fake.TemplateResponse, fake.FindError
}
