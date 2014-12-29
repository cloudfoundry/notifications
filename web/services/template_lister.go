package services

import (
	"github.com/cloudfoundry-incubator/notifications/models"
)

type TemplateSummary struct {
	Name string `json:"name"`
}

type TemplateListerInterface interface {
	List() (map[string]TemplateSummary, error)
}
type TemplateLister struct {
	templatesRepo models.TemplatesRepoInterface
	database      models.DatabaseInterface
}

func NewTemplateLister(repo models.TemplatesRepoInterface, database models.DatabaseInterface) TemplateLister {
	return TemplateLister{
		templatesRepo: repo,
		database:      database,
	}
}

func (lister TemplateLister) List() (map[string]TemplateSummary, error) {
	templates, err := lister.templatesRepo.ListIDsAndNames(lister.database.Connection())
	if err != nil {
		return map[string]TemplateSummary{}, err
	}

	templatesMap := map[string]TemplateSummary{}
	for _, template := range templates {
		if template.ID != "default" {
			templatesMap[template.ID] = TemplateSummary{Name: template.Name}
		}
	}
	return templatesMap, nil
}
