package acceptance

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/acceptance/servers"
	"github.com/cloudfoundry-incubator/notifications/config"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Get a list of all notifications", func() {
	BeforeEach(func() {
		TruncateTables()
	})

	It("allows a user to get body templates", func() {
		// Boot Fake SMTP Server
		smtpServer := servers.NewSMTP()
		smtpServer.Boot()

		// Boot Fake UAA Server
		uaaServer := servers.NewUAA()
		uaaServer.Boot()
		defer uaaServer.Close()

		// Boot Real Notifications Server
		notificationsServer := servers.NewNotifications()
		notificationsServer.Boot()
		defer notificationsServer.Close()

		// Retrieve Client UAA token
		clientID := "notifications-sender"
		env := config.NewEnvironment()
		uaaClient := uaa.NewUAA("", env.UAAHost, clientID, "secret", "")
		clientToken, err := uaaClient.GetClientToken()
		if err != nil {
			panic(err)
		}

		test := AllNotifications{notificationsServer: notificationsServer, clientToken: clientToken}
		test.SetNotifications()
		test.GetAllNotifications()
	})
})

type AllNotifications struct {
	notificationsServer servers.Notifications
	clientToken         uaa.Token
}

func (test AllNotifications) setNotifications(clientID, data string) {
	env := config.NewEnvironment()
	uaaClient := uaa.NewUAA("", env.UAAHost, clientID, "secret", "")
	clientToken, err := uaaClient.GetClientToken()
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest("PUT", test.notificationsServer.NotificationsPath(), strings.NewReader(data))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+clientToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	Expect(response.StatusCode).To(Equal(http.StatusNoContent))
}

func (test AllNotifications) SetNotifications() {
	test.setNotifications("client-123", `{
		"source_name":"source name stuff",
		"notifications":{
			"kind-asd":{
				"description":"remember stuff",
				"critical":false
			},
			"kind-abc":{
				"description":"forgot things",
				"critical":true
			}
		}
	}`)
	test.setNotifications("client-456", `{
		"source_name": "raptors",
		"notifications": {
			"dino-kind": {
				"description": "forgot things",
				"critical": true
			},
			"fossilized-kind": {
				"description": "remember stuff",
				"critical": false
			}
		}
	}`)
	test.setNotifications("client-890", `{
		"source_name": "this client has no notifications"
	}`)
}

func (test AllNotifications) GetAllNotifications() {
	request, err := http.NewRequest("GET", test.notificationsServer.NotificationsPath(), nil)
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+test.clientToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	Expect(response.StatusCode).To(Equal(http.StatusOK))
	Expect(body).To(MatchJSON(`{
		"client-123": {
			"name": "source name stuff",
			"template": "default",
			"notifications": {
				"kind-asd": {
					"description": "remember stuff",
					"template": "default",
					"critical": false
				},
				"kind-abc": {
					"description": "forgot things",
					"template": "default",
					"critical": true
				}
			}
		},
		"client-456": {
			"name": "raptors",
			"template": "default",
			"notifications": {
				"dino-kind": {
					"description": "forgot things",
					"template": "default",
					"critical": true
				},
				"fossilized-kind": {
					"description": "remember stuff",
					"template": "default",
					"critical": false
				}
			}
		},
		"client-890": {
			"name": "this client has no notifications",
			"template": "default",
			"notifications": {}
		}
	}`))
}
