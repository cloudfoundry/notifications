package services

import "github.com/cloudfoundry-incubator/notifications/models"

type TemplateUpdaterInterface interface {
    Update(models.Template) error
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

func (updater TemplateUpdater) Update(template models.Template) error {
    _, err := updater.repo.Upsert(updater.database.Connection(), template)
    if err != nil {
        return err
    }
    return nil
}
