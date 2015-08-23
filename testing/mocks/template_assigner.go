package mocks

import "github.com/cloudfoundry-incubator/notifications/db"

type TemplateAssigner struct {
	AssignToClientCall struct {
		Receives struct {
			Database   db.DatabaseInterface
			ClientID   string
			TemplateID string
		}
		Returns struct {
			Error error
		}
	}

	AssignToNotificationCall struct {
		Receives struct {
			Database       db.DatabaseInterface
			ClientID       string
			NotificationID string
			TemplateID     string
		}
		Returns struct {
			Error error
		}
	}
}

func NewTemplateAssigner() *TemplateAssigner {
	return &TemplateAssigner{}
}

func (a *TemplateAssigner) AssignToClient(database db.DatabaseInterface, clientID, templateID string) error {
	a.AssignToClientCall.Receives.Database = database
	a.AssignToClientCall.Receives.ClientID = clientID
	a.AssignToClientCall.Receives.TemplateID = templateID

	return a.AssignToClientCall.Returns.Error
}

func (a *TemplateAssigner) AssignToNotification(database db.DatabaseInterface, clientID, notificationID, templateID string) error {
	a.AssignToNotificationCall.Receives.Database = database
	a.AssignToNotificationCall.Receives.ClientID = clientID
	a.AssignToNotificationCall.Receives.NotificationID = notificationID
	a.AssignToNotificationCall.Receives.TemplateID = templateID

	return a.AssignToNotificationCall.Returns.Error
}
