package acceptance

import (
	"bytes"
	"encoding/json"
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
		env := application.NewEnvironment()
		clientID := "notifications-sender"
		uaaClient := uaa.NewUAA("", env.UAAHost, clientID, "secret", "")
		clientToken, err := uaaClient.GetClientToken()
		if err != nil {
			panic(err)
		}

		test := SendNotificationsToSpace{
			client:              support.NewClient(notificationsServer),
			clientToken:         clientToken,
			notificationsServer: notificationsServer,
			smtpServer:          smtpServer,
		}
		test.RegisterClientNotifications()
		test.CreateNewTemplate(params.Template{
			Name:    "Men in Black",
			Subject: "Aliens {{.Subject}}",
			HTML:    "<h1>Dogs</h1>{{.HTML}}",
			Text:    "Dogs\n{{.Text}}",
		})
		test.AssignTemplateToClient(clientID)
		test.SendNotificationsToSpace()
	})
})

type SendNotificationsToSpace struct {
	client              *support.Client
	TemplateID          string
	clientToken         uaa.Token
	notificationsServer servers.Notifications
	smtpServer          *servers.SMTP
}

// Make request to /registation
func (t SendNotificationsToSpace) RegisterClientNotifications() {
	code, err := t.client.Notifications.Register(t.clientToken.Access, support.RegisterClient{
		SourceName: "Notifications Sender",
		Notifications: map[string]support.RegisterNotification{
			"space-test": {
				Description: "Space Test",
			},
		},
	})

	Expect(err).NotTo(HaveOccurred())
	Expect(code).To(Equal(http.StatusNoContent))
}

func (t *SendNotificationsToSpace) CreateNewTemplate(template params.Template) {
	status, templateID, err := t.client.Templates.Create(t.clientToken.Access, template)
	Expect(err).NotTo(HaveOccurred())
	Expect(status).To(Equal(http.StatusCreated))
	Expect(templateID).NotTo(Equal(""))
	t.TemplateID = templateID
}

func (t SendNotificationsToSpace) AssignTemplateToClient(clientID string) {
	status, err := t.client.Templates.AssignToClient(t.clientToken.Access, clientID, t.TemplateID)
	Expect(err).NotTo(HaveOccurred())
	Expect(status).To(Equal(http.StatusNoContent))
}

// Make request to /spaces/:guid
func (t SendNotificationsToSpace) SendNotificationsToSpace() {
	body, err := json.Marshal(map[string]string{
		"kind_id": "space-test",
		"html":    "this is a space test",
		"text":    "this is a space test",
		"subject": "space-subject",
	})
	request, err := http.NewRequest("POST", t.notificationsServer.SpacesPath("space-123"), bytes.NewBuffer(body))
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
		return len(t.smtpServer.Deliveries)
	}, 5*time.Second).Should(Equal(1))
	delivery := t.smtpServer.Deliveries[0]

	env := application.NewEnvironment()
	Expect(delivery.Sender).To(Equal(env.Sender))
	Expect(delivery.Recipients).To(Equal([]string{"user-456@example.com"}))

	data := strings.Split(string(delivery.Data), "\n")
	Expect(data).To(ContainElement("X-CF-Client-ID: notifications-sender"))
	Expect(data).To(ContainElement("X-CF-Notification-ID: " + indexedResponses["user-456"]["notification_id"]))
	Expect(data).To(ContainElement("Subject: Aliens space-subject"))
	Expect(data).To(ContainElement("        <h1>Dogs</h1>this is a space test"))
	Expect(data).To(ContainElement("Dogs"))
	Expect(data).To(ContainElement("this is a space test"))
}
