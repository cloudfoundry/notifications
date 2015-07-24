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
	ListCall struct {
		Connection                 models.ConnectionInterface
		ReturnNotificationTypeList []models.NotificationType
		Err                        error
	}
	GetCall struct {
		Connection         models.ConnectionInterface
		notificationTypeID string
	}
	GetReturn struct {
		NotificationType models.NotificationType
		Err              error
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

func (n *NotificationTypesRepository) GetBySenderIDAndName(conn models.ConnectionInterface, senderID, name string) (models.NotificationType, error) {
	return models.NotificationType{}, nil
}

func (n *NotificationTypesRepository) List(conn models.ConnectionInterface, senderID string) ([]models.NotificationType, error) {
	n.ListCall.Connection = conn
	return n.ListCall.ReturnNotificationTypeList, n.ListCall.Err
}

func (n *NotificationTypesRepository) Get(conn models.ConnectionInterface, notificationTypeID string) (models.NotificationType, error) {
	n.GetCall.Connection = conn
	n.GetCall.notificationTypeID = notificationTypeID
	return n.GetReturn.NotificationType, n.GetReturn.Err
}
