package services

import "github.com/cloudfoundry-incubator/notifications/models"

type NotificationsUpdater struct {
	kindsRepo models.KindsRepoInterface
	database  models.DatabaseInterface
}

func NewNotificationsUpdater(kindsRepo models.KindsRepoInterface, database models.DatabaseInterface) NotificationsUpdater {
	return NotificationsUpdater{
		kindsRepo: kindsRepo,
		database:  database,
	}
}

func (updater NotificationsUpdater) Update(clientID, notificationID string, notification models.Kind) error {
	notification.ID = notificationID
	notification.ClientID = clientID

	_, err := updater.kindsRepo.Update(updater.database.Connection(), notification)
	if err != nil {
		return err
	}

	return nil
}
