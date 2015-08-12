package fakes

import (
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
)

type TemplateFinder struct {
	FindByIDCall struct {
		Receives struct {
			Database   db.DatabaseInterface
			TemplateID string
		}
		Returns struct {
			Template models.Template
			Error    error
		}
	}
}

func NewTemplateFinder() *TemplateFinder {
	return &TemplateFinder{}
}

func (tf *TemplateFinder) FindByID(database db.DatabaseInterface, templateID string) (models.Template, error) {
	tf.FindByIDCall.Receives.Database = database
	tf.FindByIDCall.Receives.TemplateID = templateID

	return tf.FindByIDCall.Returns.Template, tf.FindByIDCall.Returns.Error
}
