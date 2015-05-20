package services

import "github.com/cloudfoundry-incubator/notifications/models"

type NotificationsUpdater struct {
	kindsRepo models.KindsRepoInterface
}

func NewNotificationsUpdater(kindsRepo models.KindsRepoInterface) NotificationsUpdater {
	return NotificationsUpdater{
		kindsRepo: kindsRepo,
	}
}

func (updater NotificationsUpdater) Update(database models.DatabaseInterface, notification models.Kind) error {
	_, err := updater.kindsRepo.Update(database.Connection(), notification)
	if err != nil {
		return err
	}

	return nil
}
