package v2

import (
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/postal/common"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
)

type templateGetter interface {
	Get(connection collections.ConnectionInterface, templateID, clientID string) (collections.Template, error)
}

type TemplatesLoader struct {
	database            db.DatabaseInterface
	templatesCollection templateGetter
}

func NewTemplatesLoader(database db.DatabaseInterface, templatesCollection templateGetter) TemplatesLoader {
	return TemplatesLoader{
		database:            database,
		templatesCollection: templatesCollection,
	}
}

func (loader TemplatesLoader) LoadTemplates(clientID, kindID, templateID string) (common.Templates, error) {
	conn := loader.database.Connection()
	template, err := loader.templatesCollection.Get(conn, templateID, clientID)
	if err != nil {
		return common.Templates{}, err
	}

	return common.Templates{
		Subject: template.Subject,
		Text:    template.Text,
		HTML:    template.HTML,
	}, nil
}
