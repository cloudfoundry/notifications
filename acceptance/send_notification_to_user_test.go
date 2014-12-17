package acceptance

import (
	"net/http"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/notifications/acceptance/servers"
	"github.com/cloudfoundry-incubator/notifications/acceptance/support"
	"github.com/cloudfoundry-incubator/notifications/config"
	"github.com/cloudfoundry-incubator/notifications/web/params"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Send a notification to a user", func() {
	BeforeEach(func() {
		TruncateTables()
	})

	It("sends a single notification email to a user", func() {
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

		// Retrieve UAA token
		env := config.NewEnvironment()
		clientID := "notifications-sender"
		uaaClient := uaa.NewUAA("", env.UAAHost, clientID, "secret", "")
		clientToken, err := uaaClient.GetClientToken()
		if err != nil {
			panic(err)
		}

		createdTemplate := params.Template{
			Name:    "Star Wars",
			Subject: "Awesomeness {{.Subject}}",
			HTML:    "<p>Millenium Falcon</p>{{.HTML}}",
			Text:    "Millenium Falcon\n{{.Text}}",
		}

		t := SendNotificationToUser{
			client:              support.NewClient(notificationsServer),
			clientToken:         clientToken,
			UserID:              "user-123",
			notificationsServer: notificationsServer,
			smtpServer:          smtpServer,
		}

		t.RegisterClientNotification()
		t.CreateNewTemplate(createdTemplate)
		t.AssignTemplateToClient(clientID)
		t.SendNotificationToUser()
	})
})

type SendNotificationToUser struct {
	client              *support.Client
	TemplateID          string
	UserID              string
	clientToken         uaa.Token
	notificationsServer servers.Notifications
	smtpServer          *servers.SMTP
}

func (t SendNotificationToUser) RegisterClientNotification() {
	code, err := t.client.Notifications.Register(t.clientToken.Access, support.RegisterClient{
		SourceName: "Notifications Sender",
		Notifications: map[string]support.RegisterNotification{
			"acceptance-test": {
				Description: "Acceptance Test",
				Critical:    true,
			},
		},
	})

	Expect(err).NotTo(HaveOccurred())
	Expect(code).To(Equal(http.StatusNoContent))
}

func (t *SendNotificationToUser) CreateNewTemplate(template params.Template) {
	status, templateID, err := t.client.Templates.Create(t.clientToken.Access, template)
	Expect(err).NotTo(HaveOccurred())
	Expect(status).To(Equal(http.StatusCreated))
	Expect(templateID).NotTo(BeNil())
	t.TemplateID = templateID
}

func (t SendNotificationToUser) AssignTemplateToClient(clientID string) {
	status, err := t.client.Templates.AssignToClient(t.clientToken.Access, clientID, t.TemplateID)
	Expect(err).NotTo(HaveOccurred())
	Expect(status).To(Equal(http.StatusNoContent))
}

func (t SendNotificationToUser) SendNotificationToUser() {
	status, responses, err := t.client.Notify.User(t.clientToken.Access, t.UserID, support.Notify{
		KindID:  "acceptance-test",
		HTML:    "<p>this is an acceptance%40test</p>",
		Text:    "hello from the acceptance test",
		Subject: "my-special-subject",
	})

	Expect(err).NotTo(HaveOccurred())
	Expect(status).To(Equal(http.StatusOK))

	Expect(responses).To(HaveLen(1))
	response := responses[0]
	Expect(response.Status).To(Equal("queued"))
	Expect(response.Recipient).To(Equal(t.UserID))
	Expect(GUIDRegex.MatchString(response.NotificationID)).To(BeTrue())

	// Confirm the email message was delivered correctly
	Eventually(func() int {
		return len(t.smtpServer.Deliveries)
	}, 5*time.Second).Should(Equal(1))
	delivery := t.smtpServer.Deliveries[0]

	env := config.NewEnvironment()
	Expect(delivery.Sender).To(Equal(env.Sender))
	Expect(delivery.Recipients).To(Equal([]string{"user-123@example.com"}))

	data := strings.Split(string(delivery.Data), "\n")
	Expect(data).To(ContainElement("X-CF-Client-ID: notifications-sender"))
	Expect(data).To(ContainElement("X-CF-Notification-ID: " + response.NotificationID))
	Expect(data).To(ContainElement("Subject: Awesomeness my-special-subject"))
	Expect(data).To(ContainElement(`        <p>Millenium Falcon</p><p>this is an acceptance%40test</p>`))
	Expect(data).To(ContainElement("hello from the acceptance test"))
}
