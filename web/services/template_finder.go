package services

import "github.com/cloudfoundry-incubator/notifications/models"

type TemplateFinder struct {
	templatesRepo models.TemplatesRepoInterface
	database      models.DatabaseInterface
}

type TemplateFinderInterface interface {
	FindByID(string) (models.Template, error)
}

func NewTemplateFinder(templatesRepo models.TemplatesRepoInterface, database models.DatabaseInterface) TemplateFinder {
	return TemplateFinder{
		templatesRepo: templatesRepo,
		database:      database,
	}
}

func (finder TemplateFinder) FindByID(templateID string) (models.Template, error) {
	template, err := finder.templatesRepo.FindByID(finder.database.Connection(), templateID)
	if err != nil {
		return models.Template{}, err
	}

	return template, err
}
