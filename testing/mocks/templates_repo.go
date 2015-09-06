package mocks

import "github.com/cloudfoundry-incubator/notifications/v1/models"

type TemplatesRepo struct {
	CreateCall struct {
		Receives struct {
			Connection models.ConnectionInterface
			Template   models.Template
		}
		Returns struct {
			Template models.Template
			Error    error
		}
	}

	DestroyCall struct {
		Receives struct {
			Connection models.ConnectionInterface
			TemplateID string
		}
		Returns struct {
			Error error
		}
	}

	FindByIDCall struct {
		Receives struct {
			Connection models.ConnectionInterface
			TemplateID string
		}
		Returns struct {
			Template models.Template
			Error    error
		}
	}

	ListIDsAndNamesCall struct {
		Receives struct {
			Connection models.ConnectionInterface
		}
		Returns struct {
			Templates []models.Template
			Error     error
		}
	}

	UpdateCall struct {
		Receives struct {
			Connection models.ConnectionInterface
			TemplateID string
			Template   models.Template
		}
		Returns struct {
			Template models.Template
			Error    error
		}
	}
}

func NewTemplatesRepo() *TemplatesRepo {
	return &TemplatesRepo{}
}

func (tr *TemplatesRepo) Create(conn models.ConnectionInterface, template models.Template) (models.Template, error) {
	tr.CreateCall.Receives.Connection = conn
	tr.CreateCall.Receives.Template = template

	return tr.CreateCall.Returns.Template, tr.CreateCall.Returns.Error
}

func (tr *TemplatesRepo) Destroy(conn models.ConnectionInterface, templateID string) error {
	tr.DestroyCall.Receives.Connection = conn
	tr.DestroyCall.Receives.TemplateID = templateID

	return tr.DestroyCall.Returns.Error
}

func (tr *TemplatesRepo) FindByID(conn models.ConnectionInterface, templateID string) (models.Template, error) {
	tr.FindByIDCall.Receives.Connection = conn
	tr.FindByIDCall.Receives.TemplateID = templateID

	return tr.FindByIDCall.Returns.Template, tr.FindByIDCall.Returns.Error
}

func (tr *TemplatesRepo) ListIDsAndNames(conn models.ConnectionInterface) ([]models.Template, error) {
	tr.ListIDsAndNamesCall.Receives.Connection = conn

	return tr.ListIDsAndNamesCall.Returns.Templates, tr.ListIDsAndNamesCall.Returns.Error
}

func (tr *TemplatesRepo) Update(conn models.ConnectionInterface, templateID string, template models.Template) (models.Template, error) {
	tr.UpdateCall.Receives.Connection = conn
	tr.UpdateCall.Receives.TemplateID = templateID
	tr.UpdateCall.Receives.Template = template

	return tr.UpdateCall.Returns.Template, tr.UpdateCall.Returns.Error
}
