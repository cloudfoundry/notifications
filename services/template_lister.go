package services

import (
	"github.com/cloudfoundry-incubator/notifications/models"
)

type TemplateSummary struct {
	Name string `json:"name"`
}

type TemplateListerInterface interface {
	List(models.DatabaseInterface) (map[string]TemplateSummary, error)
}
type TemplateLister struct {
	templatesRepo templatesRepo
}

func NewTemplateLister(templatesRepo templatesRepo) TemplateLister {
	return TemplateLister{
		templatesRepo: templatesRepo,
	}
}

func (lister TemplateLister) List(database models.DatabaseInterface) (map[string]TemplateSummary, error) {
	templates, err := lister.templatesRepo.ListIDsAndNames(database.Connection())
	if err != nil {
		return map[string]TemplateSummary{}, err
	}

	templatesMap := map[string]TemplateSummary{}
	for _, template := range templates {
		if template.ID != models.DefaultTemplateID {
			templatesMap[template.ID] = TemplateSummary{Name: template.Name}
		}
	}
	return templatesMap, nil
}
