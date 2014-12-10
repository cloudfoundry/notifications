package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type TemplatesRepo struct {
	Templates                 map[string]models.Template
	FindError                 error
	UpsertError               error
	CreateError               error
	DestroyArgument           string
	DestroyError              error
	DeprecatedDestroyArgument string
	DeprecatedDestroyError    error
}

func NewTemplatesRepo() *TemplatesRepo {
	return &TemplatesRepo{
		Templates: make(map[string]models.Template),
	}
}

func (fake TemplatesRepo) FindByID(conn models.ConnectionInterface, templateID string) (models.Template, error) {
	template, ok := fake.Templates[templateID]
	if ok {
		return template, fake.FindError
	}
	return models.Template{}, models.ErrRecordNotFound{}
}

func (fake TemplatesRepo) Find(conn models.ConnectionInterface, templateName string) (models.Template, error) {
	template, ok := fake.Templates[templateName]
	if ok {
		return template, fake.FindError
	}
	return models.Template{}, models.ErrRecordNotFound{}
}

func (fake TemplatesRepo) Update(conn models.ConnectionInterface, templateID string, template models.Template) (models.Template, error) {
	fake.Templates[template.ID] = template
	return template, fake.UpsertError
}

func (fake TemplatesRepo) Upsert(conn models.ConnectionInterface, template models.Template) (models.Template, error) {
	fake.Templates[template.Name] = template
	return template, fake.UpsertError
}

func (fake *TemplatesRepo) Destroy(conn models.ConnectionInterface, templateID string) error {
	fake.DestroyArgument = templateID
	return fake.DestroyError
}

func (fake *TemplatesRepo) DeprecatedDestroy(conn models.ConnectionInterface, templateName string) error {
	fake.DeprecatedDestroyArgument = templateName
	return fake.DeprecatedDestroyError
}

func (fake *TemplatesRepo) Create(conn models.ConnectionInterface, template models.Template) (models.Template, error) {
	fake.Templates[template.Name] = template
	return template, fake.CreateError
}
