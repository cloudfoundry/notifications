package postal

import "github.com/cloudfoundry-incubator/notifications/v1/models"

type ClientsRepo interface {
	Find(connection models.ConnectionInterface, clientID string) (models.Client, error)
}

type KindsRepo interface {
	Find(connection models.ConnectionInterface, kindID string, clientID string) (models.Kind, error)
}

type MessagesRepo interface {
	Upsert(connection models.ConnectionInterface, message models.Message) (models.Message, error)
}

type ReceiptsRepo interface {
	CreateReceipts(connection models.ConnectionInterface, userGUIDs []string, clientID string, kindID string) error
}

type TemplatesRepo interface {
	FindByID(connection models.ConnectionInterface, templateID string) (models.Template, error)
}

type UnsubscribesRepo interface {
	Get(connection models.ConnectionInterface, userGUID string, clientID string, kindID string) (bool, error)
}

type GlobalUnsubscribesRepo interface {
	Get(connection models.ConnectionInterface, userGUID string) (bool, error)
}
