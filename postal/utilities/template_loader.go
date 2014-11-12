package utilities

import (
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/web/services"
)

type TemplatesLoaderInterface interface {
	LoadTemplates(string, string, string, string) (postal.Templates, error)
}

type TemplatesLoader struct {
	finder services.TemplateFinderInterface
}

func NewTemplatesLoader(finder services.TemplateFinderInterface) TemplatesLoader {
	return TemplatesLoader{
		finder: finder,
	}
}

func (loader TemplatesLoader) LoadTemplates(subjectSuffix, contentSuffix, client, kind string) (postal.Templates, error) {
	contentPath := client + "." + kind + "." + contentSuffix
	contentTemplate, err := loader.finder.Find(contentPath)
	if err != nil {
		return postal.Templates{}, err
	}

	subjectPath := client + "." + kind + "." + subjectSuffix
	subjectTemplate, err := loader.finder.Find(subjectPath)
	if err != nil {
		return postal.Templates{}, err
	}

	templates := postal.Templates{
		Subject: subjectTemplate.Text,
		Text:    contentTemplate.Text,
		HTML:    contentTemplate.HTML,
	}

	return templates, nil
}
