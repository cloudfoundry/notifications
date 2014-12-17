package acceptance

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
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

var _ = Describe("Sending notifications to all users in an organization", func() {
	BeforeEach(func() {
		TruncateTables()
	})

	It("sends a notification to each user in an organization", func() {
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
		clientID := "notifications-sender"
		uaaClient := uaa.NewUAA("", env.UAAHost, clientID, "secret", "")
		clientToken, err := uaaClient.GetClientToken()
		if err != nil {
			panic(err)
		}

		test := SendNotificationsToOrganization{
			client:      support.NewClient(notificationsServer),
			clientToken: clientToken,
		}
		test.RegisterClientNotifications(notificationsServer, clientToken)
		test.CreateNewTemplate(params.Template{
			Name:    "Gravity",
			Subject: "Coca cola {{.Subject}}",
			HTML:    "<h1>Rat</h1>{{.HTML}}",
			Text:    "Rat\n{{.Text}}",
		})
		test.AssignTemplateToClient(clientID)
		test.SendNotificationsToOrganization(notificationsServer, clientToken, smtpServer)
	})
})

type SendNotificationsToOrganization struct {
	client      *support.Client
	clientToken uaa.Token
	TemplateID  string
}

// Make request to /registation
func (t SendNotificationsToOrganization) RegisterClientNotifications(notificationsServer servers.Notifications, clientToken uaa.Token) {
	code, err := t.client.Notifications.Register(clientToken.Access, support.RegisterClient{
		SourceName: "Notifications Sender",
		Notifications: map[string]support.RegisterNotification{
			"organization-test": {
				Description: "Organization Test",
			},
		},
	})

	Expect(err).NotTo(HaveOccurred())
	Expect(code).To(Equal(http.StatusNoContent))
}

func (t *SendNotificationsToOrganization) CreateNewTemplate(template params.Template) {
	status, templateID, err := t.client.Templates.Create(t.clientToken.Access, template)
	Expect(err).NotTo(HaveOccurred())
	Expect(status).To(Equal(http.StatusCreated))
	Expect(templateID).NotTo(Equal(""))
	t.TemplateID = templateID
}

func (t SendNotificationsToOrganization) AssignTemplateToClient(clientID string) {
	status, err := t.client.Templates.AssignToClient(t.clientToken.Access, clientID, t.TemplateID)
	Expect(err).NotTo(HaveOccurred())
	Expect(status).To(Equal(http.StatusNoContent))
}

// Make request to /organization/:guid
func (t SendNotificationsToOrganization) SendNotificationsToOrganization(notificationsServer servers.Notifications, clientToken uaa.Token, smtpServer *servers.SMTP) {
	body, err := json.Marshal(map[string]string{
		"kind_id": "organization-test",
		"html":    "this is an organization test",
		"text":    "this is an organization test",
		"subject": "organization-subject",
	})
	request, err := http.NewRequest("POST", notificationsServer.OrganizationsPath("org-123"), bytes.NewBuffer(body))
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
	Expect(data).To(ContainElement("Subject: Coca cola organization-subject"))
	Expect(data).To(ContainElement("        <h1>Rat</h1>this is an organization test"))
	Expect(data).To(ContainElement("Rat"))
	Expect(data).To(ContainElement("this is an organization test"))
}
