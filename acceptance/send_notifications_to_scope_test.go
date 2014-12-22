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
		env := application.NewEnvironment()
		clientID := "notifications-sender"
		uaaClient := uaa.NewUAA("", env.UAAHost, clientID, "secret", "")
		clientToken, err := uaaClient.GetClientToken()
		if err != nil {
			panic(err)
		}

		test := SendNotificationsToUsersWithScope{
			client:              support.NewClient(notificationsServer),
			clientToken:         clientToken,
			notificationsServer: notificationsServer,
			smtpServer:          smtpServer,
			scope:               "this.scope",
		}
		test.RegisterClientNotifications()
		test.CreateNewTemplate(params.Template{
			Name:    "Frozen",
			Subject: "Food {{.Subject}}",
			HTML:    "<h1>Fish</h1>{{.HTML}}",
			Text:    "Fish\n{{.Text}}",
		})
		test.AssignTemplateToClient(clientID)
		test.SendNotificationsToScope()
	})
})

type SendNotificationsToUsersWithScope struct {
	client              *support.Client
	clientToken         uaa.Token
	TemplateID          string
	notificationsServer servers.Notifications
	smtpServer          *servers.SMTP
	scope               string
}

// Make request to /registation
func (t SendNotificationsToUsersWithScope) RegisterClientNotifications() {
	code, err := t.client.Notifications.Register(t.clientToken.Access, support.RegisterClient{
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

func (t *SendNotificationsToUsersWithScope) CreateNewTemplate(template params.Template) {
	status, templateID, err := t.client.Templates.Create(t.clientToken.Access, template)
	Expect(err).NotTo(HaveOccurred())
	Expect(status).To(Equal(http.StatusCreated))
	Expect(templateID).NotTo(Equal(""))
	t.TemplateID = templateID
}

func (t SendNotificationsToUsersWithScope) AssignTemplateToClient(clientID string) {
	status, err := t.client.Templates.AssignToClient(t.clientToken.Access, clientID, t.TemplateID)
	Expect(err).NotTo(HaveOccurred())
	Expect(status).To(Equal(http.StatusNoContent))
}

// Make request to /uaa_scopes/:scope for a scope
func (t SendNotificationsToUsersWithScope) SendNotificationsToScope() {
	t.smtpServer.Reset()

	body, err := json.Marshal(map[string]string{
		"kind_id": "scope-test",
		"html":    "this is a scope test",
		"text":    "this is a scope test",
		"subject": "scope-subject",
	})
	request, err := http.NewRequest("POST", t.notificationsServer.ScopesPath(t.scope), bytes.NewBuffer(body))
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
		return len(t.smtpServer.Deliveries)
	}, 5*time.Second).Should(Equal(1))
	delivery := t.smtpServer.Deliveries[0]

	env := application.NewEnvironment()
	Expect(delivery.Sender).To(Equal(env.Sender))
	Expect(delivery.Recipients).To(Equal([]string{"user-369@example.com"}))

	data := strings.Split(string(delivery.Data), "\n")
	Expect(data).To(ContainElement("X-CF-Client-ID: notifications-sender"))
	Expect(data).To(ContainElement("X-CF-Notification-ID: " + indexedResponses["user-369"]["notification_id"]))
	Expect(data).To(ContainElement("Subject: Food scope-subject"))
	Expect(data).To(ContainElement("        <h1>Fish</h1>this is a scope test"))
	Expect(data).To(ContainElement("this is a scope test"))
	Expect(data).To(ContainElement("Fish"))
}
