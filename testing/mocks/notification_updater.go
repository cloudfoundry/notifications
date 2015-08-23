package mocks

import (
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
)

type NotificationUpdater struct {
	UpdateCall struct {
		Receives struct {
			Database     services.DatabaseInterface
			Notification models.Kind
		}
		Returns struct {
			Error error
		}
	}
}

func (f *NotificationUpdater) Update(database services.DatabaseInterface, notification models.Kind) error {
	f.UpdateCall.Receives.Database = database
	f.UpdateCall.Receives.Notification = notification

	return f.UpdateCall.Returns.Error
}
