package acceptance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/notifications/acceptance/servers"
	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Send a notification to user with overridden template", func() {
	It("send a notification to user", func() {
		// Retrieve UAA token
		env := application.NewEnvironment()
		uaaClient := uaa.NewUAA("", env.UAAHost, "notifications-sender", "secret", "")
		clientToken, err := uaaClient.GetClientToken()
		if err != nil {
			panic(err)
		}

		t := SendOverriddenNotificationToUser{
			notificationsServer: Servers.Notifications,
			smtpServer:          Servers.SMTP,
			clientToken:         clientToken,
			textTemplate:        "text",
			htmlTemplate:        "<p>html</p>",
		}
		t.OverrideClientUserTemplate()
		t.SendNotificationToUser()
	})
})

type SendOverriddenNotificationToUser struct {
	notificationsServer servers.Notifications
	smtpServer          *servers.SMTP
	clientToken         uaa.Token
	textTemplate        string
	htmlTemplate        string
}

func (t SendOverriddenNotificationToUser) OverrideClientUserTemplate() {
	jsonBody := []byte(fmt.Sprintf(`{"text":"%s", "html":"%s"}`, t.textTemplate, t.htmlTemplate))
	request, err := http.NewRequest("PUT", t.notificationsServer.DeprecatedTemplatePath("notifications-sender."+models.UserBodyTemplateName), bytes.NewBuffer(jsonBody))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+t.clientToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	// Confirm response status code is a 204
	Expect(response.StatusCode).To(Equal(http.StatusNoContent))
}

func (t SendOverriddenNotificationToUser) SendNotificationToUser() {

	body, err := json.Marshal(map[string]string{
		"kind_id": "acceptance-test",
		"html":    "<p>this is an acceptance%40test</p>",
		"text":    "the acceptance text",
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

	Expect(data).To(ContainElement(t.textTemplate))
	Expect(data).To(ContainElement("        <p>html</p>"))
}
