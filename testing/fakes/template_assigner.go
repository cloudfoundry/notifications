package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type TemplateAssigner struct {
	AssignToClientArguments       []interface{}
	AssignToNotificationArguments []interface{}
	AssignToClientError           error
	AssignToNotificationError     error
}

func NewTemplateAssigner() *TemplateAssigner {
	return &TemplateAssigner{}
}

func (assigner *TemplateAssigner) AssignToClient(database models.DatabaseInterface, clientID, templateID string) error {
	assigner.AssignToClientArguments = []interface{}{database, clientID, templateID}
	return assigner.AssignToClientError
}

func (assigner *TemplateAssigner) AssignToNotification(database models.DatabaseInterface, clientID, notificationID, templateID string) error {
	assigner.AssignToNotificationArguments = []interface{}{database, clientID, notificationID, templateID}
	return assigner.AssignToNotificationError
}
