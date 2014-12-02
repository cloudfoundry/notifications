package acceptance

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/cloudfoundry-incubator/notifications/acceptance/servers"
	"github.com/cloudfoundry-incubator/notifications/config"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Get a list of all notifications", func() {
	BeforeEach(func() {
		TruncateTables()

		env := config.NewEnvironment()
		migrationsPath := path.Join(env.RootPath, env.ModelMigrationsDir)
		database := models.NewDatabase(env.DatabaseURL, migrationsPath)

		//notificationData
		firstClient := models.Client{
			ID:          "client-123",
			Description: "source name stuff",
		}

		firstKind := models.Kind{
			ID:          "kind-asd",
			ClientID:    firstClient.ID,
			Description: "remember stuff",
		}
		secondKind := models.Kind{
			ID:          "kind-abc",
			ClientID:    firstClient.ID,
			Description: "forgot things",
			Critical:    true,
		}

		secondClient := models.Client{
			ID:          "client-456",
			Description: "raptors",
		}

		thirdKind := models.Kind{
			ID:          "dino-kind",
			ClientID:    secondClient.ID,
			Description: "forgot things",
			Critical:    true,
		}
		fourthKind := models.Kind{
			ID:          "fossilized-kind",
			ClientID:    secondClient.ID,
			Description: "remember stuff",
		}

		thirdClient := models.Client{
			ID:          "client-890",
			Description: "this client has no notifications",
		}

		database.Connection().Insert(&firstClient)
		database.Connection().Insert(&firstKind)
		database.Connection().Insert(&secondKind)

		database.Connection().Insert(&secondClient)
		database.Connection().Insert(&thirdKind)
		database.Connection().Insert(&fourthKind)

		database.Connection().Insert(&thirdClient)
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

		test := AllNotifications{}
		test.GetAllNotifications(notificationsServer, clientToken)
	})
})

type AllNotifications struct{}

func (test AllNotifications) GetAllNotifications(notificationsServer servers.Notifications, clientToken uaa.Token) {
	request, err := http.NewRequest("GET", notificationsServer.NotificationsPath(), nil)
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+clientToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var expectedData interface{}
	json.Unmarshal([]byte(`{
		"client-123" : {
			"name":"source name stuff",
			"notifications": {
				"kind-asd": {
					"description": "remember stuff",
					"critical": false
				},
				"kind-abc" : {
					"description": "forgot things",
					"critical": true
				}
			}
		},
		"client-456" : {
			"name": "raptors",
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
		},
		"client-890" : {
			"name" : "this client has no notifications",
			"notifications": {}
		}
	}`), &expectedData)

	Expect(response.StatusCode).To(Equal(http.StatusOK))
	var actualData interface{}
	json.Unmarshal(body, &actualData)
	Expect(actualData).To(Equal(expectedData))

}
