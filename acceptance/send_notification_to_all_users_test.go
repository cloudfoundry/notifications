package acceptance

import (
	"net/http"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/notifications/acceptance/servers"
	"github.com/cloudfoundry-incubator/notifications/acceptance/support"
	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/cloudfoundry-incubator/notifications/web/params"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Send a notification to all users of UAA", func() {
	BeforeEach(func() {
		TruncateTables()
	})

	It("sends an email notification to all users of UAA", func() {
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
		clientID := "notifications-sender"
		env := application.NewEnvironment()
		uaaClient := uaa.NewUAA("", env.UAAHost, clientID, "secret", "")
		clientToken, err := uaaClient.GetClientToken()
		if err != nil {
			panic(err)
		}

		t := SendNotificationToAllUsers{
			client:              support.NewClient(notificationsServer),
			clientToken:         clientToken,
			notificationsServer: notificationsServer,
			smtpServer:          smtpServer,
		}

		t.RegisterClientNotification()
		t.CreateNewTemplate(params.Template{
			Name:    "Jurassic Park",
			Subject: "Genetics {{.Subject}}",
			HTML:    "<h1>T-Rex</h1>{{.HTML}}",
			Text:    "T-Rex\n{{.Text}}",
		})
		t.AssignTemplateToClient(clientID)
		t.SendNotificationToAllUsers()
	})
})

type SendNotificationToAllUsers struct {
	client              *support.Client
	clientToken         uaa.Token
	notificationsServer servers.Notifications
	smtpServer          *servers.SMTP
	TemplateID          string
}

// Make request to /registation
func (t SendNotificationToAllUsers) RegisterClientNotification() {
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

func (t *SendNotificationToAllUsers) CreateNewTemplate(template params.Template) {
	status, templateID, err := t.client.Templates.Create(t.clientToken.Access, template)
	Expect(err).NotTo(HaveOccurred())
	Expect(status).To(Equal(http.StatusCreated))
	Expect(templateID).NotTo(Equal(""))
	t.TemplateID = templateID
}

func (t SendNotificationToAllUsers) AssignTemplateToClient(clientID string) {
	status, err := t.client.Templates.AssignToClient(t.clientToken.Access, clientID, t.TemplateID)
	Expect(err).NotTo(HaveOccurred())
	Expect(status).To(Equal(http.StatusNoContent))
}

func (t SendNotificationToAllUsers) SendNotificationToAllUsers() {
	status, responses, err := t.client.Notify.AllUsers(t.clientToken.Access, support.Notify{
		KindID:  "acceptance-test",
		HTML:    "<p>this is an acceptance-test</p>",
		Text:    "oh no!",
		Subject: "gone awry",
	})

	Expect(err).NotTo(HaveOccurred())
	Expect(status).To(Equal(http.StatusOK))

	Expect(responses).To(HaveLen(2))

	indexedResponses := map[string]support.NotifyResponse{}
	for _, resp := range responses {
		indexedResponses[resp.Recipient] = resp
	}

	responseItem := indexedResponses["091b6583-0933-4d17-a5b6-66e54666c88e"]
	Expect(responseItem.Recipient).To(Equal("091b6583-0933-4d17-a5b6-66e54666c88e"))
	Expect(responseItem.Status).To(Equal("queued"))
	Expect(GUIDRegex.MatchString(responseItem.NotificationID)).To(BeTrue())

	responseItem = indexedResponses["943e6076-b1a5-4404-811b-a1ee9253bf56"]
	Expect(responseItem.Recipient).To(Equal("943e6076-b1a5-4404-811b-a1ee9253bf56"))
	Expect(responseItem.Status).To(Equal("queued"))
	Expect(GUIDRegex.MatchString(responseItem.NotificationID)).To(BeTrue())

	Eventually(func() int {
		return len(t.smtpServer.Deliveries)
	}, 5*time.Second).Should(Equal(2))

	recipients := []string{t.smtpServer.Deliveries[0].Recipients[0], t.smtpServer.Deliveries[1].Recipients[0]}
	Expect(recipients).To(ConsistOf([]string{"why-email@example.com", "slayer@example.com"}))

	var recipientIndex int
	if t.smtpServer.Deliveries[0].Recipients[0] == "why-email@example.com" {
		recipientIndex = 0
	} else {
		recipientIndex = 1
	}

	delivery := t.smtpServer.Deliveries[recipientIndex]
	env := application.NewEnvironment()
	Expect(delivery.Sender).To(Equal(env.Sender))

	data := strings.Split(string(delivery.Data), "\n")
	Expect(data).To(ContainElement("X-CF-Client-ID: notifications-sender"))
	Expect(data).To(ContainElement("X-CF-Notification-ID: " + indexedResponses["091b6583-0933-4d17-a5b6-66e54666c88e"].NotificationID))
	Expect(data).To(ContainElement("Subject: Genetics gone awry"))
	Expect(data).To(ContainElement("        <h1>T-Rex</h1><p>this is an acceptance-test</p>"))
	Expect(data).To(ContainElement("T-Rex"))
	Expect(data).To(ContainElement("oh no!"))
}
