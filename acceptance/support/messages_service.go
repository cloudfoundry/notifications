package support

import "encoding/json"

type MessagesService struct {
	client *Client
}

func (s MessagesService) Get(token, messageGUID string) (int, Message, error) {
	var message Message

	status, body, err := s.client.makeRequest("GET", s.client.MessagePath(messageGUID), nil, token)
	if err != nil {
		return status, message, err
	}

	err = json.NewDecoder(body).Decode(&message)
	return status, message, err
}
