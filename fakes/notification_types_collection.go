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
		Err                    error
	}
	ListCall struct {
		ReturnNotificationTypeList []collections.NotificationType
	}
}

func NewNotificationTypesCollection() *NotificationTypesCollection {
	return &NotificationTypesCollection{}
}

func (c *NotificationTypesCollection) Add(conn models.ConnectionInterface, notificationType collections.NotificationType) (collections.NotificationType, error) {
	c.AddCall.Conn = conn
	c.AddCall.NotificationType = notificationType
	return c.AddCall.ReturnNotificationType, c.AddCall.Err
}

func (c *NotificationTypesCollection) List(conn models.ConnectionInterface, senderID, clientID string) ([]collections.NotificationType, error) {
	return c.ListCall.ReturnNotificationTypeList, nil
}
