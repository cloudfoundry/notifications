package mocks

import "github.com/cloudfoundry-incubator/notifications/v1/models"

type TemplatesRepo struct {
	Templates     map[string]models.Template
	TemplatesList []models.Template
	FindError     error
	CreateError   error
	UpdateError   error
	ListError     error
	DestroyError  error

	FindByIDCall struct {
		Receives struct {
			Connection models.ConnectionInterface
			TemplateID string
		}
	}

	DestroyCall struct {
		Receives struct {
			Connection models.ConnectionInterface
			TemplateID string
		}
	}

	ListIDsAndNamesCall struct {
		Receives struct {
			Connection models.ConnectionInterface
		}
	}

	UpdateCall struct {
		Receives struct {
			Connection models.ConnectionInterface
			TemplateID string
			Template   models.Template
		}
	}
}

func NewTemplatesRepo() *TemplatesRepo {
	return &TemplatesRepo{
		Templates: make(map[string]models.Template),
	}
}

func (fake *TemplatesRepo) FindByID(conn models.ConnectionInterface, templateID string) (models.Template, error) {
	fake.FindByIDCall.Receives.Connection = conn
	fake.FindByIDCall.Receives.TemplateID = templateID

	template, ok := fake.Templates[templateID]
	if ok {
		return template, fake.FindError
	}
	return models.Template{}, models.NewRecordNotFoundError("Template %q could not be found", templateID)
}

func (fake *TemplatesRepo) Update(conn models.ConnectionInterface, templateID string, template models.Template) (models.Template, error) {
	fake.UpdateCall.Receives.Connection = conn
	fake.UpdateCall.Receives.TemplateID = templateID
	fake.UpdateCall.Receives.Template = template

	fake.Templates[template.ID] = template
	return template, fake.UpdateError
}

func (fake *TemplatesRepo) ListIDsAndNames(conn models.ConnectionInterface) ([]models.Template, error) {
	fake.ListIDsAndNamesCall.Receives.Connection = conn

	return fake.TemplatesList, fake.ListError
}

func (fake *TemplatesRepo) Destroy(conn models.ConnectionInterface, templateID string) error {
	fake.DestroyCall.Receives.Connection = conn
	fake.DestroyCall.Receives.TemplateID = templateID

	delete(fake.Templates, templateID)
	return fake.DestroyError
}

func (fake *TemplatesRepo) Create(conn models.ConnectionInterface, template models.Template) (models.Template, error) {
	fake.Templates[template.ID] = template
	return template, fake.CreateError
}
