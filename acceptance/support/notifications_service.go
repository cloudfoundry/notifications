package support

import (
	"bytes"
	"encoding/json"
)

type NotificationsService struct {
	client *Client
}

type RegisterClient struct {
	SourceName    string                          `json:"source_name"`
	Notifications map[string]RegisterNotification `json:"notifications,omitempty"`
}

type RegisterNotification struct {
	Description string `json:"description"`
	Critical    bool   `json:"critical"`
}

type NotificationsList map[string]NotificationClient

type NotificationClient struct {
	Name          string                  `json:"name"`
	Template      string                  `json:"template"`
	Notifications map[string]Notification `json:"notifications"`
}

type Notification struct {
	Description string `json:"description"`
	Template    string `json:"template"`
	Critical    bool   `json:"critical"`
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
