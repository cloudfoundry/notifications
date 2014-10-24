package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type FakeTemplatesRepo struct {
    Templates map[string]models.Template
    FindError error
}

func NewFakeTemplatesRepo() *FakeTemplatesRepo {
    return &FakeTemplatesRepo{
        Templates: make(map[string]models.Template),
    }
}

func (fake FakeTemplatesRepo) Find(conn models.ConnectionInterface, templateName string) (models.Template, error) {
    template, ok := fake.Templates[templateName]
    if ok {
        return template, fake.FindError
    }
    return models.Template{}, fake.FindError
}
