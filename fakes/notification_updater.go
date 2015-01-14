package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type NotificationUpdater struct {
	Notification models.Kind
	Error        error
}

func (f *NotificationUpdater) Update(notification models.Kind) error {
	f.Notification = notification

	return f.Error
}
