package services

import "github.com/cloudfoundry-incubator/notifications/v1/models"

type TemplateCreator struct {
	templatesRepo TemplatesRepo
}

func NewTemplateCreator(templatesRepo TemplatesRepo) TemplateCreator {
	return TemplateCreator{
		templatesRepo: templatesRepo,
	}
}

func (creator TemplateCreator) Create(database DatabaseInterface, template models.Template) (string, error) {
	newTemplate, err := creator.templatesRepo.Create(database.Connection(), template)
	if err != nil {
		return "", err
	}

	return newTemplate.ID, nil
}
