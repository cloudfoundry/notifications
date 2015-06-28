package services

import "github.com/cloudfoundry-incubator/notifications/models"

type TemplateUpdaterInterface interface {
	Update(models.DatabaseInterface, string, models.Template) error
}

type TemplateUpdater struct {
	templatesRepo TemplatesRepo
}

func NewTemplateUpdater(templatesRepo TemplatesRepo) TemplateUpdater {
	return TemplateUpdater{
		templatesRepo: templatesRepo,
	}
}

func (updater TemplateUpdater) Update(database models.DatabaseInterface, templateID string, template models.Template) error {
	_, err := updater.templatesRepo.Update(database.Connection(), templateID, template)
	if err != nil {
		return err
	}
	return nil
}
