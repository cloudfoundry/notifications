package acceptance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/notifications/acceptance/servers"
	"github.com/cloudfoundry-incubator/notifications/acceptance/support"
	"github.com/cloudfoundry-incubator/notifications/config"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sending notifications to users with certain scopes", func() {
	BeforeEach(func() {
		TruncateTables()
	})

	It("sends a notification to each user with the scope", func() {
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

		scope := "this.scope"
		test := SendNotificationsToUsersWithScope{
			client: support.NewClient(notificationsServer),
		}
		test.RegisterClientNotifications(notificationsServer, clientToken)
		test.SendNotificationsToScope(notificationsServer, clientToken, smtpServer, scope)
	})
})

type SendNotificationsToUsersWithScope struct {
	client *support.Client
}

// Make request to /registation
func (t SendNotificationsToUsersWithScope) RegisterClientNotifications(notificationsServer servers.Notifications, clientToken uaa.Token) {
	code, err := t.client.Notifications.Register(clientToken.Access, support.RegisterClient{
		SourceName: "Notifications Sender",
		Notifications: map[string]support.RegisterNotification{
			"scope-test": {
				Description: "Scope Test",
			},
		},
	})

	Expect(err).NotTo(HaveOccurred())
	Expect(code).To(Equal(http.StatusNoContent))
}

// Make request to /uaa_scopes/:scope for a scope
func (t SendNotificationsToUsersWithScope) SendNotificationsToScope(notificationsServer servers.Notifications, clientToken uaa.Token, smtpServer *servers.SMTP, scope string) {
	smtpServer.Reset()

	body, err := json.Marshal(map[string]string{
		"kind_id": "scope-test",
		"text":    "this is a scope test",
		"subject": "scope-subject",
	})
	request, err := http.NewRequest("POST", notificationsServer.ScopesPath(scope), bytes.NewBuffer(body))
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

	Expect(len(responseJSON)).To(Equal(1))

	indexedResponses := map[string]map[string]string{}
	for _, resp := range responseJSON {
		indexedResponses[resp["recipient"]] = resp
	}

	responseItem := indexedResponses["user-369"]
	Expect(responseItem["recipient"]).To(Equal("user-369"))
	Expect(responseItem["status"]).To(Equal("queued"))
	Expect(GUIDRegex.MatchString(responseItem["notification_id"])).To(BeTrue())

	// Confirm the email message was delivered correctly
	Eventually(func() int {
		return len(smtpServer.Deliveries)
	}, 5*time.Second).Should(Equal(1))
	delivery := smtpServer.Deliveries[0]

	env := config.NewEnvironment()
	Expect(delivery.Sender).To(Equal(env.Sender))
	Expect(delivery.Recipients).To(Equal([]string{"user-369@example.com"}))

	data := strings.Split(string(delivery.Data), "\n")
	Expect(data).To(ContainElement("X-CF-Client-ID: notifications-sender"))
	Expect(data).To(ContainElement("X-CF-Notification-ID: " + indexedResponses["user-369"]["notification_id"]))
	Expect(data).To(ContainElement("Subject: CF Notification: scope-subject"))
	Expect(data).To(ContainElement(`The following "Scope Test" notification was sent to you by the "Notifications Sender"`))
	Expect(data).To(ContainElement(fmt.Sprintf(`component of Cloud Foundry because you have the "%s" scope:`, scope)))
	Expect(data).To(ContainElement("this is a scope test"))
}
