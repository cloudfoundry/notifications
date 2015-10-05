package mocks

import "github.com/cloudfoundry-incubator/notifications/v1/collections"

type TemplateAssigner struct {
	AssignToClientCall struct {
		Receives struct {
			Connection collections.ConnectionInterface
			ClientID   string
			TemplateID string
		}
		Returns struct {
			Error error
		}
	}

	AssignToNotificationCall struct {
		Receives struct {
			Connection     collections.ConnectionInterface
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

func (a *TemplateAssigner) AssignToClient(connection collections.ConnectionInterface, clientID, templateID string) error {
	a.AssignToClientCall.Receives.Connection = connection
	a.AssignToClientCall.Receives.ClientID = clientID
	a.AssignToClientCall.Receives.TemplateID = templateID

	return a.AssignToClientCall.Returns.Error
}

func (a *TemplateAssigner) AssignToNotification(connection collections.ConnectionInterface, clientID, notificationID, templateID string) error {
	a.AssignToNotificationCall.Receives.Connection = connection
	a.AssignToNotificationCall.Receives.ClientID = clientID
	a.AssignToNotificationCall.Receives.NotificationID = notificationID
	a.AssignToNotificationCall.Receives.TemplateID = templateID

	return a.AssignToNotificationCall.Returns.Error
}
