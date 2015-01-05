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
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sending notifications to users with certain roles in an organization", func() {
	It("sends a notification to each user in an organization with that role", func() {
		clientID := "notifications-sender"
		clientToken := GetClientTokenFor(clientID)

		test := SendNotificationsToOrganizationRole{
			client:              support.NewClient(Servers.Notifications),
			clientToken:         clientToken,
			notificationsServer: Servers.Notifications,
			smtpServer:          Servers.SMTP,
		}
		test.RegisterClientNotifications()
		test.CreateNewTemplate(support.Template{
			Name:    "ET",
			Subject: "Phone home {{.Subject}}",
			HTML:    "<h1>Cat</h1>{{.HTML}}",
			Text:    "Cat\n{{.Text}}",
		})
		test.AssignTemplateToClient(clientID)
		test.SendNotificationsToOrganizationManagers()
		test.SendNotificationsToOrganizationAuditors()
		test.SendNotificationsToOrganizationBillingManagers()
		test.SendNotificationsToOrganizationInvalidRole()
	})
})

type SendNotificationsToOrganizationRole struct {
	client              *support.Client
	clientToken         uaa.Token
	TemplateID          string
	notificationsServer servers.Notifications
	smtpServer          *servers.SMTP
}

// Make request to /registation
func (t SendNotificationsToOrganizationRole) RegisterClientNotifications() {
	code, err := t.client.Notifications.Register(t.clientToken.Access, support.RegisterClient{
		SourceName: "Notifications Sender",
		Notifications: map[string]support.RegisterNotification{
			"organization-role-test": {
				Description: "Organization Role Test",
			},
		},
	})

	Expect(err).NotTo(HaveOccurred())
	Expect(code).To(Equal(http.StatusNoContent))
}

func (t *SendNotificationsToOrganizationRole) CreateNewTemplate(template support.Template) {
	status, templateID, err := t.client.Templates.Create(t.clientToken.Access, template)
	Expect(err).NotTo(HaveOccurred())
	Expect(status).To(Equal(http.StatusCreated))
	Expect(templateID).NotTo(Equal(""))
	t.TemplateID = templateID
}

func (t SendNotificationsToOrganizationRole) AssignTemplateToClient(clientID string) {
	status, err := t.client.Templates.AssignToClient(t.clientToken.Access, clientID, t.TemplateID)
	Expect(err).NotTo(HaveOccurred())
	Expect(status).To(Equal(http.StatusNoContent))
}

func (t SendNotificationsToOrganizationRole) SendNotificationsToOrganizationManagers() {
	t.smtpServer.Reset()

	body, err := json.Marshal(map[string]string{
		"kind_id": "organization-role-test",
		"html":    "this is another organization role test",
		"text":    "this is an organization role test",
		"subject": "organization-role-subject",
		"role":    "OrgManager",
	})
	request, err := http.NewRequest("POST", t.notificationsServer.OrganizationsPath("org-123"), bytes.NewBuffer(body))
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

	responseItem := indexedResponses["user-456"]
	Expect(responseItem["recipient"]).To(Equal("user-456"))
	Expect(responseItem["status"]).To(Equal("queued"))
	Expect(GUIDRegex.MatchString(responseItem["notification_id"])).To(BeTrue())

	// Confirm the email message was delivered correctly
	Eventually(func() int {
		return len(t.smtpServer.Deliveries)
	}, 1*time.Second).Should(Equal(1))
	delivery := t.smtpServer.Deliveries[0]

	env := application.NewEnvironment()
	Expect(delivery.Sender).To(Equal(env.Sender))
	Expect(delivery.Recipients).To(Equal([]string{"user-456@example.com"}))

	data := strings.Split(string(delivery.Data), "\n")
	Expect(data).To(ContainElement("X-CF-Client-ID: notifications-sender"))
	Expect(data).To(ContainElement("X-CF-Notification-ID: " + indexedResponses["user-456"]["notification_id"]))
	Expect(data).To(ContainElement("Subject: Phone home organization-role-subject"))
	Expect(data).To(ContainElement("Cat"))
	Expect(data).To(ContainElement("this is an organization role test"))
	Expect(data).To(ContainElement("        <h1>Cat</h1>this is another organization role test"))
}

// Make request to /organization/:guid for auditors
func (t SendNotificationsToOrganizationRole) SendNotificationsToOrganizationAuditors() {
	t.smtpServer.Reset()

	body, err := json.Marshal(map[string]string{
		"kind_id": "organization-role-test",
		"html":    "this is another organization role test",
		"text":    "this is an organization role test",
		"subject": "organization-role-subject",
		"role":    "OrgAuditor",
	})
	request, err := http.NewRequest("POST", t.notificationsServer.OrganizationsPath("org-123"), bytes.NewBuffer(body))
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

	responseItem := indexedResponses["user-123"]
	Expect(responseItem["recipient"]).To(Equal("user-123"))
	Expect(responseItem["status"]).To(Equal("queued"))
	Expect(GUIDRegex.MatchString(responseItem["notification_id"])).To(BeTrue())

	// Confirm the email message was delivered correctly
	Eventually(func() int {
		return len(t.smtpServer.Deliveries)
	}, 1*time.Second).Should(Equal(1))
	delivery := t.smtpServer.Deliveries[0]

	env := application.NewEnvironment()
	Expect(delivery.Sender).To(Equal(env.Sender))
	Expect(delivery.Recipients).To(Equal([]string{"user-123@example.com"}))

	data := strings.Split(string(delivery.Data), "\n")
	Expect(data).To(ContainElement("X-CF-Client-ID: notifications-sender"))
	Expect(data).To(ContainElement("X-CF-Notification-ID: " + indexedResponses["user-123"]["notification_id"]))
	Expect(data).To(ContainElement("Subject: Phone home organization-role-subject"))
	Expect(data).To(ContainElement("Cat"))
	Expect(data).To(ContainElement("this is an organization role test"))
	Expect(data).To(ContainElement("        <h1>Cat</h1>this is another organization role test"))
}

// Make request to /organization/:guid for billing managers
func (t SendNotificationsToOrganizationRole) SendNotificationsToOrganizationBillingManagers() {
	t.smtpServer.Reset()

	body, err := json.Marshal(map[string]string{
		"kind_id": "organization-role-test",
		"html":    "this is another organization role test",
		"text":    "this is an organization role test",
		"subject": "organization-role-subject",
		"role":    "BillingManager",
	})
	request, err := http.NewRequest("POST", t.notificationsServer.OrganizationsPath("org-123"), bytes.NewBuffer(body))
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

	responseItem := indexedResponses["user-111"]
	Expect(responseItem["recipient"]).To(Equal("user-111"))
	Expect(responseItem["status"]).To(Equal("queued"))
	Expect(GUIDRegex.MatchString(responseItem["notification_id"])).To(BeTrue())

	// Confirm the email message was delivered correctly
	Eventually(func() int {
		return len(t.smtpServer.Deliveries)
	}, 1*time.Second).Should(Equal(1))
	delivery := t.smtpServer.Deliveries[0]

	env := application.NewEnvironment()
	Expect(delivery.Sender).To(Equal(env.Sender))
	Expect(delivery.Recipients).To(Equal([]string{"user-111@example.com"}))

	data := strings.Split(string(delivery.Data), "\n")
	Expect(data).To(ContainElement("X-CF-Client-ID: notifications-sender"))
	Expect(data).To(ContainElement("X-CF-Notification-ID: " + indexedResponses["user-111"]["notification_id"]))
	Expect(data).To(ContainElement("Subject: Phone home organization-role-subject"))
	Expect(data).To(ContainElement("Cat"))
	Expect(data).To(ContainElement("this is an organization role test"))
	Expect(data).To(ContainElement("        <h1>Cat</h1>this is another organization role test"))
}

// Make request to /organization/:guid for invalid role
func (t SendNotificationsToOrganizationRole) SendNotificationsToOrganizationInvalidRole() {
	t.smtpServer.Reset()

	body, err := json.Marshal(map[string]string{
		"kind_id": "organization-role-test",
		"html":    "this is another organization role test",
		"text":    "this is an organization role test",
		"subject": "organization-role-subject",
		"role":    "bad-role",
	})
	request, err := http.NewRequest("POST", t.notificationsServer.OrganizationsPath("org-123"), bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+t.clientToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	// Confirm the request response looks correct
	Expect(response.StatusCode).To(Equal(422))

	responseJSON := map[string][]string{}
	err = json.NewDecoder(response.Body).Decode(&responseJSON)
	if err != nil {
		panic(err)
	}

	Expect(responseJSON).To(Equal(map[string][]string{
		"errors": {`"role" must be "OrgManager", "OrgAuditor", "BillingManager" or unset`},
	}))
}
