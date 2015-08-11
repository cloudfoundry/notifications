package services

import (
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/models"
)

type TemplateCreatorInterface interface {
	Create(db.DatabaseInterface, models.Template) (string, error)
}

type TemplateCreator struct {
	templatesRepo TemplatesRepo
}

func NewTemplateCreator(templatesRepo TemplatesRepo) TemplateCreator {
	return TemplateCreator{
		templatesRepo: templatesRepo,
	}
}

func (creator TemplateCreator) Create(database db.DatabaseInterface, template models.Template) (string, error) {
	newTemplate, err := creator.templatesRepo.Create(database.Connection(), template)
	if err != nil {
		return "", err
	}

	return newTemplate.ID, nil
}
