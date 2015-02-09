package support

import (
	"bytes"
	"encoding/json"
)

type NotificationsService struct {
	client *Client
}

func (s NotificationsService) Register(token string, clientToRegister RegisterClient) (int, error) {
	content, err := json.Marshal(clientToRegister)
	if err != nil {
		return 0, err
	}

	status, _, err := s.client.makeRequest("PUT", s.client.NotificationsPath(), bytes.NewBuffer(content), token)
	if err != nil {
		return 0, err
	}

	return status, nil
}

func (s NotificationsService) List(token string) (int, NotificationsList, error) {
	var list NotificationsList

	status, body, err := s.client.makeRequest("GET", s.client.NotificationsPath(), nil, token)
	if err != nil {
		return 0, list, err
	}

	err = json.NewDecoder(body).Decode(&list)
	if err != nil {
		return 0, list, err
	}

	return status, list, nil
}

func (s NotificationsService) Update(token, clientID, notificationID string, notification Notification) (int, error) {
	content, err := json.Marshal(notification)
	if err != nil {
		return 0, err
	}

	status, _, err := s.client.makeRequest("PUT", s.client.NotificationsUpdatePath(clientID, notificationID), bytes.NewBuffer(content), token)
	if err != nil {
		return 0, err
	}

	return status, nil
}
