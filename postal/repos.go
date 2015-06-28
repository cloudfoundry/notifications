package postal

import "github.com/cloudfoundry-incubator/notifications/models"

type messagesRepo interface {
	Upsert(models.ConnectionInterface, models.Message) (models.Message, error)
}

type kindsRepo interface {
	Find(models.ConnectionInterface, string, string) (models.Kind, error)
}
