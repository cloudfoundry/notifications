package postal

import (
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
)

type ClientsRepo interface {
	Find(connection db.ConnectionInterface, clientID string) (models.Client, error)
}

type KindsRepo interface {
	Find(connection db.ConnectionInterface, kindID string, clientID string) (models.Kind, error)
}

type MessagesRepo interface {
	Upsert(connection db.ConnectionInterface, message models.Message) (models.Message, error)
}

type ReceiptsRepo interface {
	CreateReceipts(connection db.ConnectionInterface, userGUIDs []string, clientID string, kindID string) error
}

type TemplatesRepo interface {
	FindByID(connection db.ConnectionInterface, templateID string) (models.Template, error)
}

type UnsubscribesRepo interface {
	Get(connection db.ConnectionInterface, userGUID string, clientID string, kindID string) (bool, error)
}

type GlobalUnsubscribesRepo interface {
	Get(connection db.ConnectionInterface, userGUID string) (bool, error)
}
