package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type FakeTemplateUpdater struct {
    UpdateArgument models.Template
    UpdateError    error
}

func NewFakeTemplateUpdater() *FakeTemplateUpdater {
    return &FakeTemplateUpdater{}
}

func (fake *FakeTemplateUpdater) Update(template models.Template) error {
    fake.UpdateArgument = template
    return fake.UpdateError
}
