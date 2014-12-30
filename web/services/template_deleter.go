package services

import "github.com/cloudfoundry-incubator/notifications/models"

type TemplateDeleterInterface interface {
	Delete(string) error
}

type TemplateDeleter struct {
	TemplatesRepo models.TemplatesRepoInterface
	Database      models.DatabaseInterface
}

func NewTemplateDeleter(repo models.TemplatesRepoInterface, database models.DatabaseInterface) TemplateDeleter {
	return TemplateDeleter{
		TemplatesRepo: repo,
		Database:      database,
	}
}

func (deleter TemplateDeleter) Delete(templateID string) error {
	connection := deleter.Database.Connection()
	return deleter.TemplatesRepo.Destroy(connection, templateID)
}
