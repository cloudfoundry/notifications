package support

import "encoding/json"

type MessagesService struct {
	client *Client
}

func (m MessagesService) Get(token, messageGUID string) (int, Message, error) {
	var message Message

	status, body, err := m.client.makeRequest("GET", m.client.server.StatusPath(messageGUID), nil, token)
	if err != nil {
		return status, message, err
	}

	err = json.NewDecoder(body).Decode(&message)
	return status, message, err
}
