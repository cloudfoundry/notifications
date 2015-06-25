package services

import "github.com/cloudfoundry-incubator/notifications/models"

type TemplateDeleterInterface interface {
	Delete(models.DatabaseInterface, string) error
}

type TemplateDeleter struct {
	templatesRepo models.TemplatesRepoInterface
}

func NewTemplateDeleter(repo models.TemplatesRepoInterface) TemplateDeleter {
	return TemplateDeleter{
		templatesRepo: repo,
	}
}

func (deleter TemplateDeleter) Delete(database models.DatabaseInterface, templateID string) error {
	connection := database.Connection()
	return deleter.templatesRepo.Destroy(connection, templateID)
}
