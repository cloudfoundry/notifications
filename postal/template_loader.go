package postal

import (
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
)

type TemplatesLoaderInterface interface {
	LoadTemplates(string, string, string) (Templates, error)
}

type TemplatesLoader struct {
	database db.DatabaseInterface

	clientsRepo         ClientsRepo
	kindsRepo           KindsRepo
	templatesRepo       TemplatesRepo
	templatesCollection templateGetter
}

type templateGetter interface {
	Get(connection collections.ConnectionInterface, templateID, clientID string) (collections.Template, error)
}

func NewTemplatesLoader(database db.DatabaseInterface, clientsRepo ClientsRepo, kindsRepo KindsRepo, templatesRepo TemplatesRepo, templatesCollection templateGetter) TemplatesLoader {
	return TemplatesLoader{
		database:            database,
		clientsRepo:         clientsRepo,
		kindsRepo:           kindsRepo,
		templatesRepo:       templatesRepo,
		templatesCollection: templatesCollection,
	}
}

func (loader TemplatesLoader) LoadTemplates(clientID, kindID, templateID string) (Templates, error) {
	conn := loader.database.Connection()

	if templateID != "" {
		template, err := loader.templatesCollection.Get(conn, templateID, clientID)
		if err != nil {
			panic(err)
		}
		return Templates{
			Subject: template.Subject,
			Text:    template.Text,
			HTML:    template.HTML,
		}, nil
	}

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

func (loader TemplatesLoader) loadTemplate(conn db.ConnectionInterface, templateID string) (Templates, error) {
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
