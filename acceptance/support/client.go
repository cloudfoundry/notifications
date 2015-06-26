package support

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type Client struct {
	host          string
	trace         bool
	hasNoRouter   bool
	Notifications *NotificationsService
	Templates     *TemplatesService
	Notify        *NotifyService
	Preferences   *PreferencesService
	Messages      *MessagesService
	HTTPClient    *http.Client
}

func NewClient(host string) *Client {
	client := &Client{
		host:        host,
		trace:       os.Getenv("TRACE") != "",
		hasNoRouter: os.Getenv("HAS_NO_ROUTER") != "",
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
	client.Messages = &MessagesService{
		client: client,
	}
	client.HTTPClient = http.DefaultClient

	return client
}

func (c Client) makeRequest(method, path string, content io.Reader, token string) (int, []byte, error) {
	request, err := http.NewRequest(method, path, content)
	if err != nil {
		return 0, []byte{}, err
	}
	c.printRequest(request)

	request.Header.Set("Authorization", "Bearer "+token)
	if c.hasNoRouter {
		request.Header.Set("X-Vcap-Request-Id", "some-totally-fake-vcap-request-id")
	}

	response, err := c.HTTPClient.Do(request)
	if err != nil {
		return 0, []byte{}, err
	}
	defer response.Body.Close()

	c.printResponse(response)

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, []byte{}, err
	}

	return response.StatusCode, responseBody, nil
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

func (c Client) SpacesPath(space string) string {
	return c.host + "/spaces/" + space
}

func (c Client) OrganizationsPath(organization string) string {
	return c.host + "/organizations/" + organization
}

func (c Client) EveryonePath() string {
	return c.host + "/everyone"
}

func (c Client) ScopesPath(scope string) string {
	return c.host + "/uaa_scopes/" + scope
}

func (c Client) UsersPath(user string) string {
	return c.host + "/users/" + user
}

func (c Client) EmailPath() string {
	return c.host + "/emails"
}

func (c Client) NotificationsPath() string {
	return c.host + "/notifications"
}

func (c Client) NotificationsUpdatePath(clientID, notificationID string) string {
	return c.host + "/clients/" + clientID + "/notifications/" + notificationID
}

func (c Client) UserPreferencesPath() string {
	return c.host + "/user_preferences"
}

func (c Client) SpecificUserPreferencesPath(userGUID string) string {
	return c.host + "/user_preferences/" + userGUID
}

func (c Client) DefaultTemplatePath() string {
	return c.host + "/default_template"
}

func (c Client) TemplatesPath() string {
	return c.host + "/templates"
}

func (c Client) TemplatePath(templateID string) string {
	return c.host + "/templates/" + templateID
}

func (c Client) TemplateAssociationsPath(templateID string) string {
	return c.host + "/templates/" + templateID + "/associations"
}

func (c Client) ClientsTemplatePath(clientID string) string {
	return c.host + "/clients/" + clientID + "/template"
}

func (c Client) ClientsNotificationsTemplatePath(clientID, notificationID string) string {
	return c.host + "/clients/" + clientID + "/notifications/" + notificationID + "/template"
}

func (c Client) MessagePath(messageID string) string {
	return c.host + "/messages/" + messageID
}
