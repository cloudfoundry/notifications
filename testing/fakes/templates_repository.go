package fakes

import (
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/v2/models"
)

type TemplatesRepository struct {
	InsertCall struct {
		Receives struct {
			Conn     db.ConnectionInterface
			Template models.Template
		}
		Returns struct {
			Template models.Template
			Err      error
		}
	}

	GetCall struct {
		Receives struct {
			Conn       db.ConnectionInterface
			TemplateID string
		}
		Returns struct {
			Template models.Template
			Err      error
		}
	}
}

func NewTemplatesRepository() *TemplatesRepository {
	return &TemplatesRepository{}
}

func (r *TemplatesRepository) Insert(conn db.ConnectionInterface, template models.Template) (models.Template, error) {
	r.InsertCall.Receives.Conn = conn
	r.InsertCall.Receives.Template = template

	return r.InsertCall.Returns.Template, r.InsertCall.Returns.Err
}

func (r *TemplatesRepository) Get(conn db.ConnectionInterface, templateID string) (models.Template, error) {
	r.GetCall.Receives.Conn = conn
	r.GetCall.Receives.TemplateID = templateID

	return r.GetCall.Returns.Template, r.GetCall.Returns.Err
}
