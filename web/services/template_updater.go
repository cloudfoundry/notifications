package services

import "github.com/cloudfoundry-incubator/notifications/models"

type TemplateUpdaterInterface interface {
	Update(string, models.Template) error
	DeprecatedUpdate(models.Template) error
}

type TemplateUpdater struct {
	repo     models.TemplatesRepoInterface
	database models.DatabaseInterface
}

func NewTemplateUpdater(repo models.TemplatesRepoInterface, database models.DatabaseInterface) TemplateUpdater {
	return TemplateUpdater{
		repo:     repo,
		database: database,
	}
}

func (updater TemplateUpdater) Update(templateID string, template models.Template) error {
	_, err := updater.repo.Update(updater.database.Connection(), templateID, template)
	if err != nil {
		return err
	}
	return nil
}

func (updater TemplateUpdater) DeprecatedUpdate(template models.Template) error {
	_, err := updater.repo.Upsert(updater.database.Connection(), template)
	if err != nil {
		return err
	}
	return nil
}
