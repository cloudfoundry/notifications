package services

import "github.com/cloudfoundry-incubator/notifications/v1/models"

type ClientsRepo interface {
	Find(connection models.ConnectionInterface, clientID string) (models.Client, error)
	FindAll(connection models.ConnectionInterface) ([]models.Client, error)
	FindAllByTemplateID(connection models.ConnectionInterface, templateID string) ([]models.Client, error)
	Update(connection models.ConnectionInterface, client models.Client) (models.Client, error)
	Upsert(connection models.ConnectionInterface, client models.Client) (models.Client, error)
}

type KindsRepo interface {
	Find(connection models.ConnectionInterface, kindID string, clientID string) (models.Kind, error)
	FindAll(connection models.ConnectionInterface) ([]models.Kind, error)
	FindAllByTemplateID(connection models.ConnectionInterface, templateID string) ([]models.Kind, error)
	Trim(connection models.ConnectionInterface, clientID string, kindIDs []string) (int, error)
	Update(connection models.ConnectionInterface, kind models.Kind) (models.Kind, error)
	Upsert(connection models.ConnectionInterface, kind models.Kind) (models.Kind, error)
}

type PreferencesRepo interface {
	FindNonCriticalPreferences(connection models.ConnectionInterface, userGUID string) ([]models.Preference, error)
}

type TemplatesRepo interface {
	Create(connection models.ConnectionInterface, template models.Template) (models.Template, error)
	Destroy(connection models.ConnectionInterface, templateID string) error
	FindByID(connection models.ConnectionInterface, templateID string) (models.Template, error)
	ListIDsAndNames(connection models.ConnectionInterface) ([]models.Template, error)
	Update(connection models.ConnectionInterface, templateID string, template models.Template) (models.Template, error)
}

type UnsubscribesRepo interface {
	Set(connection models.ConnectionInterface, userID string, clientID string, kindID string, unsubscribe bool) error
}

type GlobalUnsubscribesRepo interface {
	Get(connection models.ConnectionInterface, userGUID string) (bool, error)
	Set(connection models.ConnectionInterface, userGUID string, unsubscribe bool) error
}
