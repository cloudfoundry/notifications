package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type NotificationUpdater struct {
	UpdateCall struct {
		Arguments []interface{}
		Error     error
	}
}

func (f *NotificationUpdater) Update(database models.DatabaseInterface, notification models.Kind) error {
	f.UpdateCall.Arguments = []interface{}{database, notification}

	return f.UpdateCall.Error
}
