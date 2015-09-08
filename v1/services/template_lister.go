package services

import "github.com/cloudfoundry-incubator/notifications/v1/models"

type TemplateSummary struct {
	Name string `json:"name"`
}

type TemplateLister struct {
	templatesRepo TemplatesRepo
}

func NewTemplateLister(templatesRepo TemplatesRepo) TemplateLister {
	return TemplateLister{
		templatesRepo: templatesRepo,
	}
}

func (lister TemplateLister) List(database DatabaseInterface) (map[string]TemplateSummary, error) {
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
