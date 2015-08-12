package fakes

import (
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
)

type NotificationUpdater struct {
	UpdateCall struct {
		Receives struct {
			Database     db.DatabaseInterface
			Notification models.Kind
		}
		Returns struct {
			Error error
		}
	}
}

func (f *NotificationUpdater) Update(database db.DatabaseInterface, notification models.Kind) error {
	f.UpdateCall.Receives.Database = database
	f.UpdateCall.Receives.Notification = notification

	return f.UpdateCall.Returns.Error
}
