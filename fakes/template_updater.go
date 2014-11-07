package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type TemplateUpdater struct {
    UpdateArgument models.Template
    UpdateError    error
}

func NewTemplateUpdater() *TemplateUpdater {
    return &TemplateUpdater{}
}

func (fake *TemplateUpdater) Update(template models.Template) error {
    fake.UpdateArgument = template
    return fake.UpdateError
}
