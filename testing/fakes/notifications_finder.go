package fakes

import (
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/models"
)

type NotificationsFinder struct {
	Clients map[string]models.Client
	Kinds   map[string]models.Kind

	AllClientsAndNotificationsCall struct {
		Arguments []interface{}
		Error     error
	}

	ClientAndKindCall struct {
		Arguments []interface{}
		Error     error
	}
}

func NewNotificationsFinder() *NotificationsFinder {
	return &NotificationsFinder{
		Clients: make(map[string]models.Client),
		Kinds:   make(map[string]models.Kind),
	}
}

func (finder *NotificationsFinder) AllClientsAndNotifications(database db.DatabaseInterface) ([]models.Client, []models.Kind, error) {
	var (
		clients []models.Client
		kinds   []models.Kind
	)

	finder.AllClientsAndNotificationsCall.Arguments = []interface{}{database}

	for _, client := range finder.Clients {
		clients = append(clients, client)
	}

	for _, kind := range finder.Kinds {
		kinds = append(kinds, kind)
	}

	return clients, kinds, finder.AllClientsAndNotificationsCall.Error
}

func (finder *NotificationsFinder) ClientAndKind(database db.DatabaseInterface, clientID, kindID string) (models.Client, models.Kind, error) {
	finder.ClientAndKindCall.Arguments = []interface{}{database, clientID, kindID}
	return finder.Clients[clientID], finder.Kinds[kindID+"|"+clientID], finder.ClientAndKindCall.Error
}
