package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type NotificationUpdater struct {
	UpdateCall struct {
		Receives struct {
			Database     models.DatabaseInterface
			Notification models.Kind
		}
		Returns struct {
			Error error
		}
	}
}

func (f *NotificationUpdater) Update(database models.DatabaseInterface, notification models.Kind) error {
	f.UpdateCall.Receives.Database = database
	f.UpdateCall.Receives.Notification = notification

	return f.UpdateCall.Returns.Error
}
