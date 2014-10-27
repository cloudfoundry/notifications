package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type FakeTemplatesRepo struct {
    Templates       map[string]models.Template
    FindError       error
    UpsertError     error
    DestroyArgument string
    DestroyError    error
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

func (fake FakeTemplatesRepo) Upsert(conn models.ConnectionInterface, template models.Template) (models.Template, error) {
    fake.Templates[template.Name] = template
    return template, fake.UpsertError
}

func (fake *FakeTemplatesRepo) Destroy(conn models.ConnectionInterface, templateName string) error {
    fake.DestroyArgument = templateName
    return fake.DestroyError
}
