package services

import "github.com/cloudfoundry-incubator/notifications/models"

type TemplateCreatorInterface interface {
	Create(models.DatabaseInterface, models.Template) (string, error)
}

type TemplateCreator struct {
	repo models.TemplatesRepoInterface
}

func NewTemplateCreator(repo models.TemplatesRepoInterface) TemplateCreator {
	return TemplateCreator{
		repo: repo,
	}
}

func (creator TemplateCreator) Create(database models.DatabaseInterface, template models.Template) (string, error) {
	newTemplate, err := creator.repo.Create(database.Connection(), template)
	if err != nil {
		return "", err
	}

	return newTemplate.ID, nil
}
