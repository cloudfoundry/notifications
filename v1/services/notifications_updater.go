package services

import "github.com/cloudfoundry-incubator/notifications/v1/models"

type NotificationsUpdater struct {
	kindsRepo KindsRepo
}

func NewNotificationsUpdater(kindsRepo KindsRepo) NotificationsUpdater {
	return NotificationsUpdater{
		kindsRepo: kindsRepo,
	}
}

func (updater NotificationsUpdater) Update(database DatabaseInterface, notification models.Kind) error {
	_, err := updater.kindsRepo.Update(database.Connection(), notification)
	if err != nil {
		return err
	}

	return nil
}
