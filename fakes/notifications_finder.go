package fakes

import (
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/web/services"
)

type NotificationsFinder struct {
	Clients                     map[string]models.Client
	Kinds                       map[string]models.Kind
	ClientAndKindError          error
	ClientsWithNotifications    map[string]services.ClientWithNotifications
	AllClientNotificationsError error
}

func NewNotificationsFinder() *NotificationsFinder {
	return &NotificationsFinder{
		Clients: make(map[string]models.Client),
		Kinds:   make(map[string]models.Kind),
		ClientsWithNotifications: make(map[string]services.ClientWithNotifications),
	}
}

func (finder *NotificationsFinder) AllClientNotifications() (map[string]services.ClientWithNotifications, error) {
	return finder.ClientsWithNotifications, finder.AllClientNotificationsError
}

func (finder *NotificationsFinder) ClientAndKind(clientID, kindID string) (models.Client, models.Kind, error) {
	return finder.Clients[clientID], finder.Kinds[kindID+"|"+clientID], finder.ClientAndKindError
}
