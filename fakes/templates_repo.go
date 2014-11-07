package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type TemplatesRepo struct {
    Templates       map[string]models.Template
    FindError       error
    UpsertError     error
    DestroyArgument string
    DestroyError    error
}

func NewTemplatesRepo() *TemplatesRepo {
    return &TemplatesRepo{
        Templates: make(map[string]models.Template),
    }
}

func (fake TemplatesRepo) Find(conn models.ConnectionInterface, templateName string) (models.Template, error) {
    template, ok := fake.Templates[templateName]
    if ok {
        return template, fake.FindError
    }
    return models.Template{}, models.ErrRecordNotFound{}
}

func (fake TemplatesRepo) Upsert(conn models.ConnectionInterface, template models.Template) (models.Template, error) {
    fake.Templates[template.Name] = template
    return template, fake.UpsertError
}

func (fake *TemplatesRepo) Destroy(conn models.ConnectionInterface, templateName string) error {
    fake.DestroyArgument = templateName
    return fake.DestroyError
}
