package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type TemplatesRepo struct {
	Templates     map[string]models.Template
	TemplatesList []models.Template
	FindError     error
	CreateError   error
	UpdateError   error
	ListError     error
	DestroyError  error
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
	return models.Template{}, models.NewRecordNotFoundError("Template %q could not be found", templateID)
}

func (fake TemplatesRepo) Update(conn models.ConnectionInterface, templateID string, template models.Template) (models.Template, error) {
	fake.Templates[template.ID] = template
	return template, fake.UpdateError
}

func (fake TemplatesRepo) ListIDsAndNames(conn models.ConnectionInterface) ([]models.Template, error) {
	return fake.TemplatesList, fake.ListError
}

func (fake *TemplatesRepo) Destroy(conn models.ConnectionInterface, templateID string) error {
	delete(fake.Templates, templateID)
	return fake.DestroyError
}

func (fake *TemplatesRepo) Create(conn models.ConnectionInterface, template models.Template) (models.Template, error) {
	fake.Templates[template.ID] = template
	return template, fake.CreateError
}
