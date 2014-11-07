package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type TemplateFinder struct {
    TemplateName string
    FindError    error
    Templates    map[string]models.Template
}

func NewTemplateFinder() *TemplateFinder {
    return &TemplateFinder{
        Templates: map[string]models.Template{},
    }
}

func (fake *TemplateFinder) Find(templateName string) (models.Template, error) {
    fake.TemplateName = templateName
    return fake.Templates[templateName], fake.FindError
}
