package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type FakeTemplateFinder struct {
    TemplateName     string
    TemplateResponse models.Template
    NotificationType string
    FindError        error
}

func NewFakeTemplateFinder(templateResponse models.Template) *FakeTemplateFinder {
    return &FakeTemplateFinder{
        TemplateResponse: templateResponse,
    }
}

func (fake *FakeTemplateFinder) Find(notificationType, templateName string) (models.Template, error) {
    fake.TemplateName = templateName
    fake.NotificationType = notificationType
    return fake.TemplateResponse, fake.FindError
}
