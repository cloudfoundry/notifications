package services

import "github.com/cloudfoundry-incubator/notifications/models"

type TemplateDeleterInterface interface {
	Delete(models.DatabaseInterface, string) error
}

type TemplateDeleter struct {
	templatesRepo TemplatesRepo
}

func NewTemplateDeleter(templatesRepo TemplatesRepo) TemplateDeleter {
	return TemplateDeleter{
		templatesRepo: templatesRepo,
	}
}

func (deleter TemplateDeleter) Delete(database models.DatabaseInterface, templateID string) error {
	connection := database.Connection()
	return deleter.templatesRepo.Destroy(connection, templateID)
}
