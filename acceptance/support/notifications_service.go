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
	Notifications map[string]RegisterNotification `json:"notifications"`
}

type RegisterNotification struct {
	Description string `json:"description"`
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
