package acceptance

import (
	"net/http"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/notifications/acceptance/servers"
	"github.com/cloudfoundry-incubator/notifications/acceptance/support"
	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Send a notification to an email", func() {
	It("sends a single notification to an email", func() {
		// Retrieve UAA token
		clientID := "notifications-sender"
		env := application.NewEnvironment()
		uaaClient := uaa.NewUAA("", env.UAAHost, clientID, "secret", "")
		clientToken, err := uaaClient.GetClientToken()
		if err != nil {
			panic(err)
		}

		t := SendNotificationToEmail{
			client:              support.NewClient(Servers.Notifications),
			clientToken:         clientToken,
			notificationsServer: Servers.Notifications,
			smtpServer:          Servers.SMTP,
		}

		t.RegisterClientNotification()
		t.CreateNewTemplate(support.Template{
			Name:    "Star Trek",
			Subject: "Boldness {{.Subject}}",
			HTML:    "<p>Enterprise</p>{{.HTML}}",
			Text:    "Enterprise\n{{.Text}}",
		})
		t.AssignTemplateToClient(clientID)
		t.SendNotificationToEmail()
	})

})

type SendNotificationToEmail struct {
	client              *support.Client
	clientToken         uaa.Token
	notificationsServer servers.Notifications
	smtpServer          *servers.SMTP
	TemplateID          string
}

func (t SendNotificationToEmail) RegisterClientNotification() {
	code, err := t.client.Notifications.Register(t.clientToken.Access, support.RegisterClient{
		SourceName: "Notifications Sender",
		Notifications: map[string]support.RegisterNotification{
			"acceptance-test": {
				Description: "Acceptance Test",
			},
		},
	})

	Expect(err).NotTo(HaveOccurred())
	Expect(code).To(Equal(http.StatusNoContent))
}

func (t *SendNotificationToEmail) CreateNewTemplate(template support.Template) {
	status, templateID, err := t.client.Templates.Create(t.clientToken.Access, template)
	Expect(err).NotTo(HaveOccurred())
	Expect(status).To(Equal(http.StatusCreated))
	Expect(templateID).NotTo(Equal(""))
	t.TemplateID = templateID
}

func (t SendNotificationToEmail) AssignTemplateToClient(clientID string) {
	status, err := t.client.Templates.AssignToClient(t.clientToken.Access, clientID, t.TemplateID)
	Expect(err).NotTo(HaveOccurred())
	Expect(status).To(Equal(http.StatusNoContent))
}

func (t SendNotificationToEmail) SendNotificationToEmail() {
	status, responses, err := t.client.Notify.Email(t.clientToken.Access, support.Notify{
		KindID:  "acceptance-test",
		HTML:    "<header>this is an acceptance test</header>",
		Subject: "my-special-subject",
		To:      "John User <user@example.com>",
	})

	Expect(err).NotTo(HaveOccurred())
	Expect(status).To(Equal(http.StatusOK))

	Expect(responses).To(HaveLen(1))
	responseItem := responses[0]
	Expect(responseItem.Status).To(Equal("queued"))
	Expect(responseItem.Recipient).To(Equal("user@example.com"))
	Expect(GUIDRegex.MatchString(responseItem.NotificationID)).To(BeTrue())

	Eventually(func() int {
		return len(t.smtpServer.Deliveries)
	}, 1*time.Second).Should(Equal(1))
	delivery := t.smtpServer.Deliveries[0]

	env := application.NewEnvironment()
	Expect(delivery.Sender).To(Equal(env.Sender))
	Expect(delivery.Recipients).To(Equal([]string{"user@example.com"}))

	data := strings.Split(string(delivery.Data), "\n")
	Expect(data).To(ContainElement("X-CF-Client-ID: notifications-sender"))
	Expect(data).To(ContainElement("X-CF-Notification-ID: " + responseItem.NotificationID))
	Expect(data).To(ContainElement("Subject: Boldness my-special-subject"))
	Expect(data).To(ContainElement("        <p>Enterprise</p><header>this is an acceptance test</header>"))
}
