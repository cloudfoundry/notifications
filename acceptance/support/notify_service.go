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
	Subject           string
	HTML              string
	Text              string
	KindID            string
	ReplyTo           string
	SourceDescription string
}

func (nr notifyRequest) Merge(n Notify) notifyRequest {
	nr.Subject = n.Subject
	nr.HTML = n.HTML
	nr.Text = n.Text
	nr.KindID = n.KindID
	nr.ReplyTo = n.ReplyTo

	return nr
}

func (s NotifyService) notify(token, path string, notify Notify, reqBody notifyRequest) (int, []NotifyResponse, error) {
	var responses []NotifyResponse

	reqBody = reqBody.Merge(notify)
	body, err := json.Marshal(reqBody)
	if err != nil {
		return 0, responses, err
	}

	status, responseBody, err := s.client.makeRequest("POST", path, bytes.NewBuffer(body), token)
	if err != nil {
		return 0, responses, err
	}

	if status == http.StatusOK {
		err = json.Unmarshal(responseBody, &responses)
		if err != nil {
			return 0, responses, err
		}
	}

	return status, responses, nil
}

func (s NotifyService) User(token, userGUID string, notify Notify) (int, []NotifyResponse, error) {
	return s.notify(token, s.client.UsersPath(userGUID), notify, notifyRequest{})
}

func (s NotifyService) AllUsers(token string, notify Notify) (int, []NotifyResponse, error) {
	return s.notify(token, s.client.EveryonePath(), notify, notifyRequest{})
}

func (s NotifyService) Email(token, email string, notify Notify) (int, []NotifyResponse, error) {
	return s.notify(token, s.client.EmailPath(), notify, notifyRequest{
		To: email,
	})
}

func (s NotifyService) OrganizationRole(token, organizationGUID, role string, notify Notify) (int, []NotifyResponse, error) {
	return s.notify(token, s.client.OrganizationsPath(organizationGUID), notify, notifyRequest{
		Role: role,
	})
}

func (s NotifyService) Organization(token, organizationGUID string, notify Notify) (int, []NotifyResponse, error) {
	return s.notify(token, s.client.OrganizationsPath(organizationGUID), notify, notifyRequest{})
}

func (s NotifyService) Scope(token, scope string, notify Notify) (int, []NotifyResponse, error) {
	return s.notify(token, s.client.ScopesPath(scope), notify, notifyRequest{})
}

func (s NotifyService) Space(token, spaceGUID string, notify Notify) (int, []NotifyResponse, error) {
	return s.notify(token, s.client.SpacesPath(spaceGUID), notify, notifyRequest{})
}
