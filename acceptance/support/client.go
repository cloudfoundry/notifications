package support

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/cloudfoundry-incubator/notifications/acceptance/servers"
)

type Client struct {
	server        servers.Notifications
	trace         bool
	Notifications *NotificationsService
	Templates     *TemplatesService
	Notify        *NotifyService
	Preferences   *PreferencesService
}

func NewClient(server servers.Notifications) *Client {
	client := &Client{
		server: server,
		trace:  os.Getenv("TRACE") != "",
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
	client.Preferences = &PreferencesService{
		client: client,
	}

	return client
}

func (c Client) makeRequest(method, path string, content io.Reader, token string) (int, io.Reader, error) {
	request, err := http.NewRequest(method, path, content)
	if err != nil {
		return 0, nil, err
	}
	c.printRequest(request)

	request.Header.Set("Authorization", "Bearer "+token)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return 0, nil, err
	}
	c.printResponse(response)

	return response.StatusCode, response.Body, nil
}

func (c Client) printRequest(request *http.Request) {
	if c.trace {
		buffer := bytes.NewBuffer([]byte{})
		body := bytes.NewBuffer([]byte{})
		if request.Body != nil {
			_, err := io.Copy(io.MultiWriter(buffer, body), request.Body)
			if err != nil {
				panic(err)
			}
		}

		request.Body = ioutil.NopCloser(body)

		fmt.Printf("[REQ] %s %s %s\n", request.Method, request.URL.String(), buffer.String())
	}
}

func (c Client) printResponse(response *http.Response) {
	if c.trace {
		buffer := bytes.NewBuffer([]byte{})
		body := bytes.NewBuffer([]byte{})
		if response.Body != nil {
			_, err := io.Copy(io.MultiWriter(buffer, body), response.Body)
			if err != nil {
				panic(err)
			}
		}

		response.Body = ioutil.NopCloser(body)

		fmt.Printf("[RES] %s %s\n", response.Status, buffer.String())
	}
}
