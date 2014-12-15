package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type NotificationsFinder struct {
	Clients                         map[string]models.Client
	Kinds                           map[string]models.Kind
	ClientAndKindError              error
	AllClientsAndNotificationsError error
}

func NewNotificationsFinder() *NotificationsFinder {
	return &NotificationsFinder{
		Clients: make(map[string]models.Client),
		Kinds:   make(map[string]models.Kind),
	}
}

func (finder *NotificationsFinder) AllClientsAndNotifications() ([]models.Client, []models.Kind, error) {
	var clients []models.Client
	var kinds []models.Kind
	for _, client := range finder.Clients {
		clients = append(clients, client)
	}

	for _, kind := range finder.Kinds {
		kinds = append(kinds, kind)
	}

	return clients, kinds, finder.AllClientsAndNotificationsError
}

func (finder *NotificationsFinder) ClientAndKind(clientID, kindID string) (models.Client, models.Kind, error) {
	return finder.Clients[clientID], finder.Kinds[kindID+"|"+clientID], finder.ClientAndKindError
}
