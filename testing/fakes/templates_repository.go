package fakes

import "github.com/cloudfoundry-incubator/notifications/v2/models"

type TemplatesRepository struct {
	InsertCall struct {
		Receives struct {
			Connection models.ConnectionInterface
			Template   models.Template
		}
		Returns struct {
			Template models.Template
			Error    error
		}
	}

	GetCall struct {
		Receives struct {
			Connection models.ConnectionInterface
			TemplateID string
		}
		Returns struct {
			Template models.Template
			Error    error
		}
	}

	DeleteCall struct {
		Receives struct {
			Connection models.ConnectionInterface
			TemplateID string
		}
		Returns struct {
			Error error
		}
	}
}

func NewTemplatesRepository() *TemplatesRepository {
	return &TemplatesRepository{}
}

func (r *TemplatesRepository) Insert(conn models.ConnectionInterface, template models.Template) (models.Template, error) {
	r.InsertCall.Receives.Connection = conn
	r.InsertCall.Receives.Template = template

	return r.InsertCall.Returns.Template, r.InsertCall.Returns.Error
}

func (r *TemplatesRepository) Get(conn models.ConnectionInterface, templateID string) (models.Template, error) {
	r.GetCall.Receives.Connection = conn
	r.GetCall.Receives.TemplateID = templateID

	return r.GetCall.Returns.Template, r.GetCall.Returns.Error
}

func (r *TemplatesRepository) Delete(conn models.ConnectionInterface, templateID string) error {
	r.DeleteCall.Receives.Connection = conn
	r.DeleteCall.Receives.TemplateID = templateID

	return r.DeleteCall.Returns.Error
}
