package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type TemplateFinder struct {
	FindByIDCall struct {
		Receives struct {
			Database   models.DatabaseInterface
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

func (tf *TemplateFinder) FindByID(database models.DatabaseInterface, templateID string) (models.Template, error) {
	tf.FindByIDCall.Receives.Database = database
	tf.FindByIDCall.Receives.TemplateID = templateID

	return tf.FindByIDCall.Returns.Template, tf.FindByIDCall.Returns.Error
}
