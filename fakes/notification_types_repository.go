package fakes

import (
	"github.com/cloudfoundry-incubator/notifications/models"
)

type NotificationTypesRepository struct {
	InsertCall struct {
		Connection             models.ConnectionInterface
		NotificationType       models.NotificationType
		ReturnNotificationType models.NotificationType
		Err                    error
	}
}

func NewNotificationTypesRepository() *NotificationTypesRepository {
	return &NotificationTypesRepository{}
}

func (n *NotificationTypesRepository) Insert(conn models.ConnectionInterface, notificationType models.NotificationType) (models.NotificationType, error) {
	n.InsertCall.NotificationType = notificationType
	n.InsertCall.Connection = conn
	return n.InsertCall.ReturnNotificationType, n.InsertCall.Err
}
