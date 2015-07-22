package fakes

import (
	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/models"
)

type NotificationTypesCollection struct {
	AddCall struct {
		NotificationType       collections.NotificationType
		Conn                   models.ConnectionInterface
		ReturnNotificationType collections.NotificationType
	}
}

func NewNotificationTypesCollection() *NotificationTypesCollection {
	return &NotificationTypesCollection{}
}

func (c *NotificationTypesCollection) Add(conn models.ConnectionInterface, notificationType collections.NotificationType) (collections.NotificationType, error) {
	c.AddCall.Conn = conn
	c.AddCall.NotificationType = notificationType
	return c.AddCall.ReturnNotificationType, nil
}
