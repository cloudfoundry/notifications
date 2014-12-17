package utilities

import (
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/web/services"
)

type TemplatesLoaderInterface interface {
	DeprecatedLoadTemplates(string, string, string, string) (postal.Templates, error)
	LoadTemplates(string, string, string, string) (postal.Templates, error)
}

type TemplatesLoader struct {
	finder        services.TemplateFinderInterface
	database      models.DatabaseInterface
	clientsRepo   models.ClientsRepoInterface
	kindsRepo     models.KindsRepoInterface
	templatesRepo models.TemplatesRepoInterface
}

func NewTemplatesLoader(finder services.TemplateFinderInterface, database models.DatabaseInterface, clientsRepo models.ClientsRepoInterface,
	kindsRepo models.KindsRepoInterface, templatesRepo models.TemplatesRepoInterface) TemplatesLoader {

	return TemplatesLoader{
		finder:        finder,
		database:      database,
		clientsRepo:   clientsRepo,
		kindsRepo:     kindsRepo,
		templatesRepo: templatesRepo,
	}
}

func (loader TemplatesLoader) DeprecatedLoadTemplates(subjectSuffix, contentSuffix, client, kind string) (postal.Templates, error) {
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

func (loader TemplatesLoader) LoadTemplates(clientID, kindID, contentSuffix, subjectSuffix string) (postal.Templates, error) {
	conn := loader.database.Connection()

	kind, err := loader.kindsRepo.Find(conn, kindID, clientID)
	if err != nil {
		return postal.Templates{}, err
	}

	if kind.Template != "" {
		return loader.loadTemplate(conn, kind.Template)
	}

	client, err := loader.clientsRepo.Find(conn, clientID)
	if err != nil {
		return postal.Templates{}, err
	}

	if client.Template != "" {
		return loader.loadTemplate(conn, client.Template)
	}

	return loader.DeprecatedLoadTemplates(subjectSuffix, contentSuffix, clientID, kindID)
}

func (loader TemplatesLoader) loadTemplate(conn models.ConnectionInterface, templateID string) (postal.Templates, error) {
	template, err := loader.templatesRepo.FindByID(conn, templateID)
	if err != nil {
		return postal.Templates{}, err
	}

	return postal.Templates{
		Subject: template.Subject,
		Text:    template.Text,
		HTML:    template.HTML,
	}, nil
}
