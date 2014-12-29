package support

import (
	"io"
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
	client.Notifications = &NotificationsService{
		client: client,
	}
	client.Templates = &TemplatesService{
		client: client,
		Default: &DefaultTemplateService{
			client: client,
		},
	}
	client.Notify = &NotifyService{
		client: client,
	}

	return client
}

func (c Client) makeRequest(method, path string, content io.Reader, token string) (int, io.Reader, error) {
	request, err := http.NewRequest(method, path, content)
	if err != nil {
		return 0, nil, err
	}

	request.Header.Set("Authorization", "Bearer "+token)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return 0, nil, err
	}

	return response.StatusCode, response.Body, nil
}
