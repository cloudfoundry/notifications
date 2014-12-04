package services

import "github.com/cloudfoundry-incubator/notifications/models"

type TemplateCreatorInterface interface {
	Create(models.Template) (string, error)
}

type TemplateCreator struct {
	repo     models.TemplatesRepoInterface
	database models.DatabaseInterface
}

func NewTemplateCreator(repo models.TemplatesRepoInterface, database models.DatabaseInterface) TemplateCreator {
	return TemplateCreator{
		repo:     repo,
		database: database,
	}
}

func (creator TemplateCreator) Create(template models.Template) (string, error) {
	newTemplate, err := creator.repo.Create(creator.database.Connection(), template)
	if err != nil {
		return "", err
	}

	return newTemplate.ID, nil
}
