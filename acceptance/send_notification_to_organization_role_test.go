package acceptance

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/notifications/acceptance/servers"
	"github.com/cloudfoundry-incubator/notifications/config"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sending notifications to users with certain roles in an organization", func() {
	BeforeEach(func() {
		TruncateTables()
	})

	It("sends a notification to each user in an organization with that role", func() {
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

		test := SendNotificationsToOrganizationRole{}
		test.RegisterClientNotifications(notificationsServer, clientToken)
		test.SendNotificationsToOrganizationManagers(notificationsServer, clientToken, smtpServer)
		test.SendNotificationsToOrganizationAuditors(notificationsServer, clientToken, smtpServer)
		test.SendNotificationsToOrganizationBillingManagers(notificationsServer, clientToken, smtpServer)
		test.SendNotificationsToOrganizationInvalidRole(notificationsServer, clientToken, smtpServer)
	})
})

type SendNotificationsToOrganizationRole struct{}

// Make request to /registation
func (t SendNotificationsToOrganizationRole) RegisterClientNotifications(notificationsServer servers.Notifications, clientToken uaa.Token) {
	body, err := json.Marshal(map[string]interface{}{
		"source_description": "Notifications Sender",
		"kinds": []map[string]string{
			{
				"id":          "organization-role-test",
				"description": "Organization Role Test",
			},
		},
	})
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest("PUT", notificationsServer.RegistrationPath(), bytes.NewBuffer(body))
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

	// Confirm response status code looks ok
	Expect(response.StatusCode).To(Equal(http.StatusOK))

}

// Make request to /organization/:guid for managers
func (t SendNotificationsToOrganizationRole) SendNotificationsToOrganizationManagers(notificationsServer servers.Notifications, clientToken uaa.Token, smtpServer *servers.SMTP) {
	smtpServer.Reset()

	body, err := json.Marshal(map[string]string{
		"kind_id": "organization-role-test",
		"text":    "this is an organization role test",
		"subject": "organization-role-subject",
		"role":    "OrgManager",
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
		return len(smtpServer.Deliveries)
	}, 5*time.Second).Should(Equal(1))
	delivery := smtpServer.Deliveries[0]

	env := config.NewEnvironment()
	Expect(delivery.Sender).To(Equal(env.Sender))
	Expect(delivery.Recipients).To(Equal([]string{"user-456@example.com"}))

	data := strings.Split(string(delivery.Data), "\n")
	Expect(data).To(ContainElement("X-CF-Client-ID: notifications-sender"))
	Expect(data).To(ContainElement("X-CF-Notification-ID: " + indexedResponses["user-456"]["notification_id"]))
	Expect(data).To(ContainElement("Subject: CF Notification: organization-role-subject"))
	Expect(data).To(ContainElement(`The following "Organization Role Test" notification was sent to you by the "Notifications Sender"`))
	Expect(data).To(ContainElement(`component of Cloud Foundry because you are a member of the "notifications-service" organization:`))
	Expect(data).To(ContainElement("this is an organization role test"))
}

// Make request to /organization/:guid for auditors
func (t SendNotificationsToOrganizationRole) SendNotificationsToOrganizationAuditors(notificationsServer servers.Notifications, clientToken uaa.Token, smtpServer *servers.SMTP) {
	smtpServer.Reset()

	body, err := json.Marshal(map[string]string{
		"kind_id": "organization-role-test",
		"text":    "this is an organization role test",
		"subject": "organization-role-subject",
		"role":    "OrgAuditor",
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
		return len(smtpServer.Deliveries)
	}, 5*time.Second).Should(Equal(1))
	delivery := smtpServer.Deliveries[0]

	env := config.NewEnvironment()
	Expect(delivery.Sender).To(Equal(env.Sender))
	Expect(delivery.Recipients).To(Equal([]string{"user-123@example.com"}))

	data := strings.Split(string(delivery.Data), "\n")
	Expect(data).To(ContainElement("X-CF-Client-ID: notifications-sender"))
	Expect(data).To(ContainElement("X-CF-Notification-ID: " + indexedResponses["user-123"]["notification_id"]))
	Expect(data).To(ContainElement("Subject: CF Notification: organization-role-subject"))
	Expect(data).To(ContainElement(`The following "Organization Role Test" notification was sent to you by the "Notifications Sender"`))
	Expect(data).To(ContainElement(`component of Cloud Foundry because you are a member of the "notifications-service" organization:`))
	Expect(data).To(ContainElement("this is an organization role test"))
}

// Make request to /organization/:guid for billing managers
func (t SendNotificationsToOrganizationRole) SendNotificationsToOrganizationBillingManagers(notificationsServer servers.Notifications, clientToken uaa.Token, smtpServer *servers.SMTP) {
	smtpServer.Reset()

	body, err := json.Marshal(map[string]string{
		"kind_id": "organization-role-test",
		"text":    "this is an organization role test",
		"subject": "organization-role-subject",
		"role":    "BillingManager",
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
		return len(smtpServer.Deliveries)
	}, 5*time.Second).Should(Equal(1))
	delivery := smtpServer.Deliveries[0]

	env := config.NewEnvironment()
	Expect(delivery.Sender).To(Equal(env.Sender))
	Expect(delivery.Recipients).To(Equal([]string{"user-111@example.com"}))

	data := strings.Split(string(delivery.Data), "\n")
	Expect(data).To(ContainElement("X-CF-Client-ID: notifications-sender"))
	Expect(data).To(ContainElement("X-CF-Notification-ID: " + indexedResponses["user-111"]["notification_id"]))
	Expect(data).To(ContainElement("Subject: CF Notification: organization-role-subject"))
	Expect(data).To(ContainElement(`The following "Organization Role Test" notification was sent to you by the "Notifications Sender"`))
	Expect(data).To(ContainElement(`component of Cloud Foundry because you are a member of the "notifications-service" organization:`))
	Expect(data).To(ContainElement("this is an organization role test"))
}

// Make request to /organization/:guid for invalid role
func (t SendNotificationsToOrganizationRole) SendNotificationsToOrganizationInvalidRole(notificationsServer servers.Notifications, clientToken uaa.Token, smtpServer *servers.SMTP) {
	smtpServer.Reset()

	body, err := json.Marshal(map[string]string{
		"kind_id": "organization-role-test",
		"text":    "this is an organization role test",
		"subject": "organization-role-subject",
		"role":    "bad-role",
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
	Expect(response.StatusCode).To(Equal(422))

	responseJSON := map[string][]string{}
	err = json.Unmarshal(body, &responseJSON)
	if err != nil {
		panic(err)
	}

	Expect(responseJSON).To(Equal(map[string][]string{
		"errors": {`"role" must be "OrgManager", "OrgAuditor", "BillingManager" or unset`},
	}))
}
