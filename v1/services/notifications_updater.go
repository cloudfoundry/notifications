package services

import (
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/models"
)

type NotificationsUpdaterInterface interface {
	Update(db.DatabaseInterface, models.Kind) error
}

type NotificationsUpdater struct {
	kindsRepo KindsRepo
}

func NewNotificationsUpdater(kindsRepo KindsRepo) NotificationsUpdater {
	return NotificationsUpdater{
		kindsRepo: kindsRepo,
	}
}

func (updater NotificationsUpdater) Update(database db.DatabaseInterface, notification models.Kind) error {
	_, err := updater.kindsRepo.Update(database.Connection(), notification)
	if err != nil {
		return err
	}

	return nil
}