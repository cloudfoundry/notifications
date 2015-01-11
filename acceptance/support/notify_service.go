package support

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type NotifyService struct {
	client *Client
}

type Notify struct {
	Subject string
	HTML    string
	Text    string
	KindID  string
}

type notifyRequest struct {
	To      string `json:"to,omitempty"`
	Role    string `json:"role,omitempty"`
	Subject string `json:"subject"`
	HTML    string `json:"html"`
	Text    string `json:"text"`
	KindID  string `json:"kind_id,omitempty"`
}

func (nr notifyRequest) Merge(n Notify) notifyRequest {
	nr.Subject = n.Subject
	nr.HTML = n.HTML
	nr.Text = n.Text
	nr.KindID = n.KindID

	return nr
}

type NotifyResponse struct {
	Status         string `json:"status"`
	Recipient      string `json:"recipient"`
	NotificationID string `json:"notification_id"`
}

func (service NotifyService) notify(token, path string, notify Notify, reqBody notifyRequest) (int, []NotifyResponse, error) {
	var responses []NotifyResponse

	reqBody = reqBody.Merge(notify)
	body, err := json.Marshal(reqBody)
	if err != nil {
		return 0, responses, err
	}

	status, responseBody, err := service.client.makeRequest("POST", path, bytes.NewBuffer(body), token)
	if err != nil {
		return 0, responses, err
	}

	if status == http.StatusOK {
		err = json.NewDecoder(responseBody).Decode(&responses)
		if err != nil {
			return 0, responses, err
		}
	}

	return status, responses, nil
}

func (service NotifyService) User(token, userGUID string, notify Notify) (int, []NotifyResponse, error) {
	return service.notify(token, service.client.server.UsersPath(userGUID), notify, notifyRequest{})
}

func (service NotifyService) AllUsers(token string, notify Notify) (int, []NotifyResponse, error) {
	return service.notify(token, service.client.server.EveryonePath(), notify, notifyRequest{})
}

func (service NotifyService) Email(token, email string, notify Notify) (int, []NotifyResponse, error) {
	return service.notify(token, service.client.server.EmailPath(), notify, notifyRequest{
		To: email,
	})
}

func (service NotifyService) OrganizationRole(token, organizationGUID, role string, notify Notify) (int, []NotifyResponse, error) {
	return service.notify(token, service.client.server.OrganizationsPath(organizationGUID), notify, notifyRequest{
		Role: role,
	})
}

func (service NotifyService) Organization(token, organizationGUID string, notify Notify) (int, []NotifyResponse, error) {
	return service.notify(token, service.client.server.OrganizationsPath(organizationGUID), notify, notifyRequest{})
}

func (service NotifyService) Scope(token, scope string, notify Notify) (int, []NotifyResponse, error) {
	return service.notify(token, service.client.server.ScopesPath(scope), notify, notifyRequest{})
}
