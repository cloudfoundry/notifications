package support

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type NotificationType struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Critical    bool   `json:"critical"`
	TemplateID  string `json:"template_id"`
}

type NotificationTypesService struct {
	config Config
}

func NewNotificationTypesService(config Config) NotificationTypesService {
	return NotificationTypesService{
		config: config,
	}
}

func (n NotificationTypesService) Create(senderID, name, description, templateID string, critical bool, token string) (NotificationType, error) {
	var notificationType NotificationType

	content, err := json.Marshal(map[string]interface{}{
		"name":        name,
		"description": description,
		"critical":    critical,
		"template_id": templateID,
	})
	if err != nil {
		return notificationType, err
	}

	status, body, err := NewClient(n.config).makeRequest("POST", n.config.Host+"/senders/"+senderID+"/notification_types", bytes.NewBuffer(content), token)
	if err != nil {
		return notificationType, err
	}

	if status != http.StatusCreated {
		return notificationType, UnexpectedStatusError{status, string(body)}
	}

	err = json.Unmarshal(body, &notificationType)
	if err != nil {
		return notificationType, err
	}

	return notificationType, nil
}

func (n NotificationTypesService) Show(senderID, notificationTypeID, token string) (NotificationType, error) {
	var notificationType NotificationType

	status, body, err := NewClient(n.config).makeRequest("GET", n.config.Host+"/senders/"+senderID+"/notification_types/"+notificationTypeID, nil, token)
	if err != nil {
		return notificationType, err
	}

	if status == http.StatusNotFound {
		return notificationType, NotFoundError{status, string(body)}
	} else if status != http.StatusOK {
		return notificationType, UnexpectedStatusError{status, string(body)}
	}

	err = json.Unmarshal(body, &notificationType)
	if err != nil {
		return notificationType, err
	}

	return notificationType, nil
}
