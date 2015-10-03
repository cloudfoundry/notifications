package v1

import (
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/postal/common"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
)

type clientFinder interface {
	Find(connection models.ConnectionInterface, clientID string) (models.Client, error)
}

type kindFinder interface {
	Find(connection models.ConnectionInterface, kindID string, clientID string) (models.Kind, error)
}

type templateFinder interface {
	FindByID(connection models.ConnectionInterface, templateID string) (models.Template, error)
}

type TemplatesLoader struct {
	database db.DatabaseInterface

	clientsRepo   clientFinder
	kindsRepo     kindFinder
	templatesRepo templateFinder
}

func NewTemplatesLoader(database db.DatabaseInterface, clientsRepo clientFinder, kindsRepo kindFinder, templatesRepo templateFinder) TemplatesLoader {
	return TemplatesLoader{
		database:      database,
		clientsRepo:   clientsRepo,
		kindsRepo:     kindsRepo,
		templatesRepo: templatesRepo,
	}
}

func (loader TemplatesLoader) LoadTemplates(clientID, kindID, templateID string) (common.Templates, error) {
	conn := loader.database.Connection()

	if kindID != "" {
		kind, err := loader.kindsRepo.Find(conn, kindID, clientID)
		if err != nil {
			return common.Templates{}, err
		}

		if kind.TemplateID != models.DefaultTemplateID {
			return loader.loadTemplate(conn, kind.TemplateID)
		}
	}

	client, err := loader.clientsRepo.Find(conn, clientID)
	if err != nil {
		return common.Templates{}, err
	}

	return loader.loadTemplate(conn, client.TemplateID)
}

func (loader TemplatesLoader) loadTemplate(conn db.ConnectionInterface, templateID string) (common.Templates, error) {
	template, err := loader.templatesRepo.FindByID(conn, templateID)
	if err != nil {
		return common.Templates{}, err
	}

	return common.Templates{
		Subject: template.Subject,
		Text:    template.Text,
		HTML:    template.HTML,
	}, nil
}
