package services

import "github.com/cloudfoundry-incubator/notifications/models"

type TemplateUpdaterInterface interface {
	Update(models.DatabaseInterface, string, models.Template) error
}

type TemplateUpdater struct {
	repo models.TemplatesRepoInterface
}

func NewTemplateUpdater(repo models.TemplatesRepoInterface) TemplateUpdater {
	return TemplateUpdater{
		repo: repo,
	}
}

func (updater TemplateUpdater) Update(database models.DatabaseInterface, templateID string, template models.Template) error {
	_, err := updater.repo.Update(database.Connection(), templateID, template)
	if err != nil {
		return err
	}
	return nil
}
