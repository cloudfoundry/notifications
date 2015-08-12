package services

import (
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
)

type ClientsRepo interface {
	Find(connection db.ConnectionInterface, clientID string) (models.Client, error)
	FindAll(connection db.ConnectionInterface) ([]models.Client, error)
	FindAllByTemplateID(connection db.ConnectionInterface, templateID string) ([]models.Client, error)
	Update(connection db.ConnectionInterface, client models.Client) (models.Client, error)
	Upsert(connection db.ConnectionInterface, client models.Client) (models.Client, error)
}

type KindsRepo interface {
	Find(connection db.ConnectionInterface, kindID string, clientID string) (models.Kind, error)
	FindAll(connection db.ConnectionInterface) ([]models.Kind, error)
	FindAllByTemplateID(connection db.ConnectionInterface, templateID string) ([]models.Kind, error)
	Trim(connection db.ConnectionInterface, clientID string, kindIDs []string) (int, error)
	Update(connection db.ConnectionInterface, kind models.Kind) (models.Kind, error)
	Upsert(connection db.ConnectionInterface, kind models.Kind) (models.Kind, error)
}

type PreferencesRepo interface {
	FindNonCriticalPreferences(connection db.ConnectionInterface, userGUID string) ([]models.Preference, error)
}

type TemplatesRepo interface {
	Create(connection db.ConnectionInterface, template models.Template) (models.Template, error)
	Destroy(connection db.ConnectionInterface, templateID string) error
	FindByID(connection db.ConnectionInterface, templateID string) (models.Template, error)
	ListIDsAndNames(connection db.ConnectionInterface) ([]models.Template, error)
	Update(connection db.ConnectionInterface, templateID string, template models.Template) (models.Template, error)
}

type UnsubscribesRepo interface {
	Set(connection db.ConnectionInterface, userID string, clientID string, kindID string, unsubscribe bool) error
}

type GlobalUnsubscribesRepo interface {
	Get(connection db.ConnectionInterface, userGUID string) (bool, error)
	Set(connection db.ConnectionInterface, userGUID string, unsubscribe bool) error
}
