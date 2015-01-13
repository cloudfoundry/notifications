package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type NotificationUpdater struct {
	ClientID     string
	ID           string
	Notification models.Kind
	Error        error
}

func (f *NotificationUpdater) Update(clientID, notificationID string, notification models.Kind) error {
	f.ClientID = clientID
	f.ID = notificationID
	f.Notification = notification

	return f.Error
}
