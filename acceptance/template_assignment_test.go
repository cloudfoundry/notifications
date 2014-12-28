package acceptance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/acceptance/servers"
	"github.com/cloudfoundry-incubator/notifications/acceptance/support"
	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/cloudfoundry-incubator/notifications/web/params"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Assign Templates", func() {
	It("Creates a template and then assigns it", func() {
		// Retrieve Client UAA token
		clientID := "notifications-admin"
		env := application.NewEnvironment()
		uaaClient := uaa.NewUAA("", env.UAAHost, clientID, "secret", "")
		clientToken, err := uaaClient.GetClientToken()
		if err != nil {
			panic(err)
		}

		notificationID := "acceptance-test"
		createdTemplate := params.Template{
			Name:    "Star Wars",
			Subject: "Awesomeness",
			HTML:    "<p>Millenium Falcon</p>",
			Text:    "Millenium Falcon",
		}

		t := AssignTemplate{
			client:              support.NewClient(Servers.Notifications),
			notificationsServer: Servers.Notifications,
			clientToken:         clientToken,
		}
		t.RegisterClientNotification(notificationID)
		t.CreateNewTemplate(createdTemplate)
		t.AssignTemplateToClient(clientID)
		t.ConfirmClientTemplateAssignment(clientID)
		t.AssignTemplateToNotification(clientID, notificationID)
		t.ConfirmNotificationTemplateAssignment(clientID, notificationID)
	})
})

type NotificationsResponse map[string]struct {
	Name          string `json:"name"`
	Template      string `json:"template"`
	Notifications map[string]struct {
		Description string `json:"description"`
		Critical    bool   `json:"critical"`
		Template    string `json:"template"`
	} `json:"notifications"`
}

type AssignTemplate struct {
	client              *support.Client
	notificationsServer servers.Notifications
	clientToken         uaa.Token
	TemplateID          string
}

func (test AssignTemplate) RegisterClientNotification(notificationID string) {
	code, err := test.client.Notifications.Register(test.clientToken.Access, support.RegisterClient{
		SourceName: "Notifications Sender",
		Notifications: map[string]support.RegisterNotification{
			notificationID: {
				Description: "Acceptance Test",
				Critical:    true,
			},
		},
	})

	Expect(err).NotTo(HaveOccurred())
	Expect(code).To(Equal(http.StatusNoContent))
}

func (test *AssignTemplate) CreateNewTemplate(template params.Template) {
	status, templateID := test.createTemplateHelper(template)
	test.TemplateID = templateID
	Expect(status).To(Equal(http.StatusCreated))
	Expect(templateID).NotTo(BeNil())
}

func (test *AssignTemplate) AssignTemplateToClient(clientID string) {
	status := test.assignTemplateHelper(clientID)
	Expect(status).To(Equal(http.StatusNoContent))
}

func (test *AssignTemplate) ConfirmClientTemplateAssignment(clientID string) {
	status, notifications := test.getNotifications()
	Expect(status).To(Equal(http.StatusOK))
	Expect(notifications[clientID].Template).To(Equal(test.TemplateID))
}

func (test *AssignTemplate) AssignTemplateToNotification(clientID, notificationID string) {
	status := test.assignNotificationHelper(clientID, notificationID)
	Expect(status).To(Equal(http.StatusNoContent))
}

func (test *AssignTemplate) ConfirmNotificationTemplateAssignment(clientID, notificationID string) {
	status, notifications := test.getNotifications()
	Expect(status).To(Equal(http.StatusOK))
	clientNotifications := notifications[clientID].Notifications
	Expect(clientNotifications[notificationID].Template).To(Equal(test.TemplateID))
}

func (test *AssignTemplate) createTemplateHelper(templateToCreate params.Template) (int, string) {
	jsonBody, err := json.Marshal(templateToCreate)
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest("POST", test.notificationsServer.TemplatesBasePath(), bytes.NewBuffer(jsonBody))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+test.clientToken.Access)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	var JSON struct {
		TemplateID string `json:"template_id"`
	}

	err = json.NewDecoder(response.Body).Decode(&JSON)
	if err != nil {
		panic(err)
	}

	return response.StatusCode, JSON.TemplateID
}

func (test *AssignTemplate) assignTemplateHelper(clientID string) int {
	request, err := http.NewRequest("PUT", test.notificationsServer.ClientsTemplatePath(clientID), bytes.NewBuffer([]byte(fmt.Sprintf(`{"template":%q}`, test.TemplateID))))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+test.clientToken.Access)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	return response.StatusCode
}

func (test *AssignTemplate) assignNotificationHelper(clientID, notificationID string) int {
	request, err := http.NewRequest("PUT", test.notificationsServer.ClientsNotificationsTemplatePath(clientID, notificationID), bytes.NewBuffer([]byte(fmt.Sprintf(`{"template":%q}`, test.TemplateID))))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+test.clientToken.Access)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	return response.StatusCode
}

func (test *AssignTemplate) getNotifications() (int, NotificationsResponse) {
	request, err := http.NewRequest("GET", test.notificationsServer.NotificationsPath(), nil)
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+test.clientToken.Access)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	JSON := NotificationsResponse{}
	err = json.NewDecoder(response.Body).Decode(&JSON)
	if err != nil {
		panic(err)
	}

	return response.StatusCode, JSON
}
