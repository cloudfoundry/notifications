package services

import "github.com/cloudfoundry-incubator/notifications/db"

type TemplateDeleterInterface interface {
	Delete(db.DatabaseInterface, string) error
}

type TemplateDeleter struct {
	templatesRepo TemplatesRepo
}

func NewTemplateDeleter(templatesRepo TemplatesRepo) TemplateDeleter {
	return TemplateDeleter{
		templatesRepo: templatesRepo,
	}
}

func (deleter TemplateDeleter) Delete(database db.DatabaseInterface, templateID string) error {
	connection := database.Connection()
	return deleter.templatesRepo.Destroy(connection, templateID)
}
