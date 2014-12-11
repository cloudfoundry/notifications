package services

import (
	"github.com/cloudfoundry-incubator/notifications/models"
)

type TemplateMetadata struct {
	Name string `json:"name"`
}

type TemplateListerInterface interface {
	List() (map[string]TemplateMetadata, error)
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

func (lister TemplateLister) List() (map[string]TemplateMetadata, error) {
	templates, err := lister.templatesRepo.ListIDsAndNames(lister.database.Connection())
	if err != nil {
		return map[string]TemplateMetadata{}, err
	}

	templatesMap := map[string]TemplateMetadata{}
	for _, template := range templates {
		templatesMap[template.ID] = TemplateMetadata{Name: template.Name}
	}
	return templatesMap, nil
}
