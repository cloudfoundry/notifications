package mocks

import (
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
)

type NotificationsFinder struct {
	AllClientsAndNotificationsCall struct {
		Receives struct {
			Database db.DatabaseInterface
		}
		Returns struct {
			Clients []models.Client
			Kinds   []models.Kind
			Error   error
		}
	}

	ClientAndKindCall struct {
		Receives struct {
			Database db.DatabaseInterface
			ClientID string
			KindID   string
		}
		Returns struct {
			Client models.Client
			Kind   models.Kind
			Error  error
		}
	}
}

func NewNotificationsFinder() *NotificationsFinder {
	return &NotificationsFinder{}
}

func (f *NotificationsFinder) AllClientsAndNotifications(database db.DatabaseInterface) ([]models.Client, []models.Kind, error) {
	f.AllClientsAndNotificationsCall.Receives.Database = database

	return f.AllClientsAndNotificationsCall.Returns.Clients, f.AllClientsAndNotificationsCall.Returns.Kinds, f.AllClientsAndNotificationsCall.Returns.Error
}

func (f *NotificationsFinder) ClientAndKind(database db.DatabaseInterface, clientID, kindID string) (models.Client, models.Kind, error) {
	f.ClientAndKindCall.Receives.Database = database
	f.ClientAndKindCall.Receives.ClientID = clientID
	f.ClientAndKindCall.Receives.KindID = kindID

	return f.ClientAndKindCall.Returns.Client, f.ClientAndKindCall.Returns.Kind, f.ClientAndKindCall.Returns.Error
}
