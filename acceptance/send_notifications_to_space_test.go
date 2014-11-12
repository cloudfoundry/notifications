package acceptance

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/notifications/acceptance/servers"
	"github.com/cloudfoundry-incubator/notifications/config"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sending notifications to all users in a space", func() {
	BeforeEach(func() {
		TruncateTables()
	})

	It("sends a notification to each user in a space", func() {
		// Boot Fake SMTP Server
		smtpServer := servers.NewSMTP()
		smtpServer.Boot()

		// Boot Fake UAA Server
		uaaServer := servers.NewUAA()
		uaaServer.Boot()
		defer uaaServer.Close()

		// Boot Fake CC Server
		ccServer := servers.NewCC()
		ccServer.Boot()
		defer ccServer.Close()

		// Boot Real Notifications Server
		notificationsServer := servers.NewNotifications()
		notificationsServer.Boot()
		defer notificationsServer.Close()

		// Retrieve UAA token
		env := config.NewEnvironment()
		uaaClient := uaa.NewUAA("", env.UAAHost, "notifications-sender", "secret", "")
		clientToken, err := uaaClient.GetClientToken()
		if err != nil {
			panic(err)
		}

		test := SendNotificationsToSpace{}
		test.RegisterClientNotifications(notificationsServer, clientToken)
		test.SendNotificationsToSpace(notificationsServer, clientToken, smtpServer)
	})
})

type SendNotificationsToSpace struct{}

// Make request to /registation
func (t SendNotificationsToSpace) RegisterClientNotifications(notificationsServer servers.Notifications, clientToken uaa.Token) {
	body, err := json.Marshal(map[string]interface{}{
		"source_description": "Notifications Sender",
		"kinds": []map[string]string{
			{
				"id":          "space-test",
				"description": "Space Test",
			},
		},
	})
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest("PUT", notificationsServer.RegistrationPath(), bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+clientToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	// Confirm response status code looks ok
	Expect(response.StatusCode).To(Equal(http.StatusOK))

}

// Make request to /spaces/:guid
func (t SendNotificationsToSpace) SendNotificationsToSpace(notificationsServer servers.Notifications, clientToken uaa.Token, smtpServer *servers.SMTP) {
	body, err := json.Marshal(map[string]string{
		"kind_id": "space-test",
		"text":    "this is a space test",
		"subject": "space-subject",
	})
	request, err := http.NewRequest("POST", notificationsServer.SpacesPath("space-123"), bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+clientToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	// Confirm the request response looks correct
	Expect(response.StatusCode).To(Equal(http.StatusOK))

	responseJSON := []map[string]string{}
	err = json.Unmarshal(body, &responseJSON)
	if err != nil {
		panic(err)
	}

	Expect(len(responseJSON)).To(Equal(3))

	indexedResponses := map[string]map[string]string{}
	for _, resp := range responseJSON {
		indexedResponses[resp["recipient"]] = resp
	}

	responseItem := indexedResponses["user-456"]
	Expect(responseItem["recipient"]).To(Equal("user-456"))
	Expect(responseItem["status"]).To(Equal("queued"))
	Expect(GUIDRegex.MatchString(responseItem["notification_id"])).To(BeTrue())

	responseItem = indexedResponses["user-789"]
	Expect(responseItem["recipient"]).To(Equal("user-789"))
	Expect(responseItem["status"]).To(Equal("queued"))
	Expect(GUIDRegex.MatchString(responseItem["notification_id"])).To(BeTrue())

	responseItem = indexedResponses["user-000"]
	Expect(responseItem["recipient"]).To(Equal("user-000"))
	Expect(responseItem["status"]).To(Equal("queued"))
	Expect(GUIDRegex.MatchString(responseItem["notification_id"])).To(BeTrue())

	// Confirm the email message was delivered correctly
	Eventually(func() int {
		return len(smtpServer.Deliveries)
	}, 5*time.Second).Should(Equal(1))
	delivery := smtpServer.Deliveries[0]

	env := config.NewEnvironment()
	Expect(delivery.Sender).To(Equal(env.Sender))
	Expect(delivery.Recipients).To(Equal([]string{"user-456@example.com"}))

	data := strings.Split(string(delivery.Data), "\n")
	Expect(data).To(ContainElement("X-CF-Client-ID: notifications-sender"))
	Expect(data).To(ContainElement("X-CF-Notification-ID: " + indexedResponses["user-456"]["notification_id"]))
	Expect(data).To(ContainElement("Subject: CF Notification: space-subject"))
	Expect(data).To(ContainElement(`The following "Space Test" notification was sent to you by the "Notifications Sender"`))
	Expect(data).To(ContainElement(`component of Cloud Foundry because you are a member of the "notifications-service" space`))
	Expect(data).To(ContainElement(`in the "notifications-service" organization:`))
	Expect(data).To(ContainElement("this is a space test"))
}
