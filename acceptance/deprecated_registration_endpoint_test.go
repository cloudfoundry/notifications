package acceptance

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/notifications/acceptance/servers"
	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Send a notification to a user, using the deprecated /registration endpoint", func() {
	It("sends a single notification email to a user", func() {
		// Retrieve UAA token
		env := application.NewEnvironment()
		uaaClient := uaa.NewUAA("", env.UAAHost, "notifications-sender", "secret", "")
		clientToken, err := uaaClient.GetClientToken()
		if err != nil {
			panic(err)
		}

		t := DeprecatedRegistrationEndpoint{
			notificationsServer: Servers.Notifications,
			smtpServer:          Servers.SMTP,
			clientToken:         clientToken,
		}
		t.RegisterClientNotification()
		t.SendNotificationToUser()
	})
})

type DeprecatedRegistrationEndpoint struct {
	notificationsServer servers.Notifications
	smtpServer          *servers.SMTP
	clientToken         uaa.Token
}

// Make request to /registation
func (t DeprecatedRegistrationEndpoint) RegisterClientNotification() {
	body, err := json.Marshal(map[string]interface{}{
		"source_description": "Notifications Sender",
		"kinds": []map[string]string{
			{
				"id":          "acceptance-test",
				"description": "Acceptance Test",
			},
		},
	})
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest("PUT", t.notificationsServer.RegistrationPath(), bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+t.clientToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	// Confirm response status code looks ok
	Expect(response.StatusCode).To(Equal(http.StatusOK))
}

// Make request to /users/:guid
func (t DeprecatedRegistrationEndpoint) SendNotificationToUser() {
	body, err := json.Marshal(map[string]string{
		"kind_id": "acceptance-test",
		"html":    "<p>this is an acceptance%40test</p>",
		"subject": "my-special-subject",
	})
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest("POST", t.notificationsServer.UsersPath("user-123"), bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+t.clientToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	// Confirm the request response looks correct
	Expect(response.StatusCode).To(Equal(http.StatusOK))

	responseJSON := []map[string]string{}
	err = json.NewDecoder(response.Body).Decode(&responseJSON)
	if err != nil {
		panic(err)
	}

	Expect(len(responseJSON)).To(Equal(1))
	responseItem := responseJSON[0]
	Expect(responseItem["status"]).To(Equal("queued"))
	Expect(responseItem["recipient"]).To(Equal("user-123"))
	Expect(GUIDRegex.MatchString(responseItem["notification_id"])).To(BeTrue())

	// Confirm the email message was delivered correctly
	Eventually(func() int {
		return len(t.smtpServer.Deliveries)
	}, 5*time.Second).Should(Equal(1))
	delivery := t.smtpServer.Deliveries[0]

	env := application.NewEnvironment()
	Expect(delivery.Sender).To(Equal(env.Sender))
	Expect(delivery.Recipients).To(Equal([]string{"user-123@example.com"}))

	data := strings.Split(string(delivery.Data), "\n")
	Expect(data).To(ContainElement("X-CF-Client-ID: notifications-sender"))
	Expect(data).To(ContainElement("X-CF-Notification-ID: " + responseItem["notification_id"]))
	Expect(data).To(ContainElement("Subject: CF Notification: my-special-subject"))
	Expect(data).To(ContainElement(`<p>The following "Acceptance Test" notification was sent to you directly by the`))
	Expect(data).To(ContainElement(`    "Notifications Sender" component of Cloud Foundry:</p>`))
	Expect(data).To(ContainElement("<p>this is an acceptance%40test</p>"))
}
