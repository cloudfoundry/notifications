package services

import "github.com/cloudfoundry-incubator/notifications/models"

type NotificationsUpdaterInterface interface {
	Update(models.DatabaseInterface, models.Kind) error
}

type NotificationsUpdater struct {
	kindsRepo kindsRepo
}

func NewNotificationsUpdater(kindsRepo kindsRepo) NotificationsUpdater {
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
