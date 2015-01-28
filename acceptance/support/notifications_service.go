package support

import (
	"bytes"
	"encoding/json"
)

type NotificationsService struct {
	client *Client
}

func (n NotificationsService) Register(token string, clientToRegister RegisterClient) (int, error) {
	content, err := json.Marshal(clientToRegister)
	if err != nil {
		return 0, err
	}

	status, _, err := n.client.makeRequest("PUT", n.client.server.NotificationsPath(), bytes.NewBuffer(content), token)
	if err != nil {
		return 0, err
	}

	return status, nil
}

func (n NotificationsService) List(token string) (int, NotificationsList, error) {
	var list NotificationsList

	status, body, err := n.client.makeRequest("GET", n.client.server.NotificationsPath(), nil, token)
	if err != nil {
		return 0, list, err
	}

	err = json.NewDecoder(body).Decode(&list)
	if err != nil {
		return 0, list, err
	}

	return status, list, nil
}

func (n NotificationsService) Update(token, clientID, notificationID string, notification Notification) (int, error) {
	content, err := json.Marshal(notification)
	if err != nil {
		return 0, err
	}

	status, _, err := n.client.makeRequest("PUT", n.client.server.NotificationsUpdatePath(clientID, notificationID), bytes.NewBuffer(content), token)
	if err != nil {
		return 0, err
	}

	return status, nil
}
