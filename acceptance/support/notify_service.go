package support

import (
	"bytes"
	"encoding/json"
)

type NotifyService struct {
	client *Client
}

type Notify struct {
	KindID  string `json:"kind_id"`
	HTML    string `json:"html"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
	To      string `json:"to,omitempty"`
}

type NotifyResponse struct {
	Status         string `json:"status"`
	Recipient      string `json:"recipient"`
	NotificationID string `json:"notification_id"`
}

func (n NotifyService) User(token, userID string, notify Notify) (int, []NotifyResponse, error) {
	var status int
	var responses []NotifyResponse

	body, err := json.Marshal(notify)
	if err != nil {
		return status, responses, err
	}

	request, err := n.client.makeRequest("POST", n.client.server.UsersPath(userID), bytes.NewBuffer(body), token)
	if err != nil {
		return status, responses, err
	}

	status, body, err = n.client.do(request)
	if err != nil {
		return status, responses, err
	}

	err = json.Unmarshal(body, &responses)
	if err != nil {
		return status, responses, err
	}

	return status, responses, nil
}

func (n NotifyService) AllUsers(token string, notify Notify) (int, []NotifyResponse, error) {
	var status int
	var responses []NotifyResponse

	body, err := json.Marshal(notify)
	if err != nil {
		return status, responses, err
	}

	request, err := n.client.makeRequest("POST", n.client.server.EveryonePath(), bytes.NewBuffer(body), token)
	if err != nil {
		return status, responses, err
	}

	status, body, err = n.client.do(request)
	if err != nil {
		return status, responses, err
	}

	err = json.Unmarshal(body, &responses)
	if err != nil {
		return status, responses, err
	}

	return status, responses, nil
}

func (n NotifyService) Email(token string, notify Notify) (int, []NotifyResponse, error) {
	var status int
	var responses []NotifyResponse

	body, err := json.Marshal(notify)
	if err != nil {
		return status, responses, err
	}

	request, err := n.client.makeRequest("POST", n.client.server.EmailPath(), bytes.NewBuffer(body), token)
	if err != nil {
		return status, responses, err
	}

	status, body, err = n.client.do(request)
	if err != nil {
		return status, responses, err
	}

	err = json.Unmarshal(body, &responses)
	if err != nil {
		return status, responses, err
	}

	return status, responses, nil
}
