package postal

import (
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
)

type templateGetter interface {
	Get(connection collections.ConnectionInterface, templateID, clientID string) (collections.Template, error)
}

type V2TemplatesLoader struct {
	database            db.DatabaseInterface
	templatesCollection templateGetter
}

func NewV2TemplatesLoader(database db.DatabaseInterface, templatesCollection templateGetter) V2TemplatesLoader {
	return V2TemplatesLoader{
		database:            database,
		templatesCollection: templatesCollection,
	}
}

func (loader V2TemplatesLoader) LoadTemplates(clientID, kindID, templateID string) (Templates, error) {
	conn := loader.database.Connection()
	template, err := loader.templatesCollection.Get(conn, templateID, clientID)
	if err != nil {
		return Templates{}, err
	}

	return Templates{
		Subject: template.Subject,
		Text:    template.Text,
		HTML:    template.HTML,
	}, nil
}
