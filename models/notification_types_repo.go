package models

import (
	"database/sql"
	"strings"
)

type NotificationTypesRepository struct {
	guidGenerator guidGeneratorFunc
}

func NewNotificationTypesRepository(guidGenerator guidGeneratorFunc) NotificationTypesRepository {
	return NotificationTypesRepository{
		guidGenerator: guidGenerator,
	}
}

func (n NotificationTypesRepository) Insert(connection ConnectionInterface, notificationType NotificationType) (NotificationType, error) {
	id, err := n.guidGenerator()
	if err != nil {
		panic(err)
	}

	notificationType.ID = id.String()
	err = connection.Insert(&notificationType)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			err = DuplicateRecordError{}
		}
		return notificationType, err
	}

	return notificationType, nil
}

func (n NotificationTypesRepository) GetBySenderIDAndName(connection ConnectionInterface, senderID, name string) (NotificationType, error) {
	notificationType := NotificationType{}
	err := connection.SelectOne(&notificationType, "SELECT * FROM `notification_types` WHERE `sender_id` = ? AND `name` = ?", senderID, name)
	if err != nil {
		if err == sql.ErrNoRows {
			err = NewRecordNotFoundError("Campaign type with sender_id %q and name %q could not be found", senderID, name)
		}
		return notificationType, err
	}

	return notificationType, nil
}

func (n NotificationTypesRepository) List(connection ConnectionInterface, senderID string) ([]NotificationType, error) {
	notificationTypeList := []NotificationType{}
	_, err := connection.Select(&notificationTypeList, "SELECT * FROM `notification_types` WHERE `sender_id` = ?", senderID)
	if err != nil {
		panic(err)
	}

	return notificationTypeList, nil
}

func (n NotificationTypesRepository) Get(connection ConnectionInterface, notificationTypeID string) (NotificationType, error) {
	notificationType := NotificationType{}
	err := connection.SelectOne(&notificationType, "SELECT * FROM `notification_types` WHERE `id` = ?", notificationTypeID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = NewRecordNotFoundError("Campaign type with id %q could not be found", notificationTypeID)
		}
		return notificationType, err
	}

	return notificationType, nil
}
