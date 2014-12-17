package support

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/acceptance/servers"
)

type Client struct {
	server        servers.Notifications
	Notifications *NotificationsService
	Templates     *TemplatesService
	Notify        *NotifyService
}

func NewClient(server servers.Notifications) *Client {
	client := &Client{
		server: server,
	}
	client.Notifications = &NotificationsService{client: client}
	client.Templates = &TemplatesService{client: client}
	client.Notify = &NotifyService{client: client}

	return client
}

func (c Client) makeRequest(method, path string, content io.Reader, token string) (*http.Request, error) {
	request, err := http.NewRequest(method, path, content)
	if err != nil {
		return request, err
	}

	request.Header.Set("Authorization", "Bearer "+token)

	return request, nil
}

func (c Client) do(request *http.Request) (int, []byte, error) {
	var status int
	var body []byte

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return status, body, err
	}

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return status, body, err
	}

	return response.StatusCode, body, nil
}
