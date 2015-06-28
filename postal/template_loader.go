package postal

import (
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/services"
)

type TemplatesLoaderInterface interface {
	LoadTemplates(string, string) (Templates, error)
}

type TemplatesLoader struct {
	finder        services.TemplateFinderInterface
	database      models.DatabaseInterface
	clientsRepo   clientsRepo
	kindsRepo     kindsRepo
	templatesRepo models.TemplatesRepoInterface
}

type clientsRepo interface {
	Find(models.ConnectionInterface, string) (models.Client, error)
}

func NewTemplatesLoader(finder services.TemplateFinderInterface, database models.DatabaseInterface, clientsRepo clientsRepo, kindsRepo kindsRepo, templatesRepo models.TemplatesRepoInterface) TemplatesLoader {

	return TemplatesLoader{
		finder:        finder,
		database:      database,
		clientsRepo:   clientsRepo,
		kindsRepo:     kindsRepo,
		templatesRepo: templatesRepo,
	}
}

func (loader TemplatesLoader) LoadTemplates(clientID, kindID string) (Templates, error) {
	conn := loader.database.Connection()

	if kindID != "" {
		kind, err := loader.kindsRepo.Find(conn, kindID, clientID)
		if err != nil {
			return Templates{}, err
		}

		if kind.TemplateID != models.DefaultTemplateID {
			return loader.loadTemplate(conn, kind.TemplateID)
		}
	}

	client, err := loader.clientsRepo.Find(conn, clientID)
	if err != nil {
		return Templates{}, err
	}

	return loader.loadTemplate(conn, client.TemplateID)
}

func (loader TemplatesLoader) loadTemplate(conn models.ConnectionInterface, templateID string) (Templates, error) {
	template, err := loader.templatesRepo.FindByID(conn, templateID)
	if err != nil {
		return Templates{}, err
	}

	return Templates{
		Subject: template.Subject,
		Text:    template.Text,
		HTML:    template.HTML,
	}, nil
}
