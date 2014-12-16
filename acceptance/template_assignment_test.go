package acceptance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/acceptance/servers"
	"github.com/cloudfoundry-incubator/notifications/config"
	"github.com/cloudfoundry-incubator/notifications/web/params"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Assign Templates", func() {
	BeforeEach(func() {
		TruncateTables()
	})

	It("Creates a template and then assigns it", func() {
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

		// Retrieve Client UAA token
		clientID := "notifications-admin"
		env := config.NewEnvironment()
		uaaClient := uaa.NewUAA("", env.UAAHost, clientID, "secret", "")
		clientToken, err := uaaClient.GetClientToken()
		if err != nil {
			panic(err)
		}

		notificationID := "acceptance-test"

		test := AssignTemplate{
			notificationsServer: notificationsServer,
			clientToken:         clientToken,
		}

		createdTemplate := params.Template{
			Name:    "Star Wars",
			Subject: "Awesomeness",
			HTML:    "<p>Millenium Falcon</p>",
			Text:    "Millenium Falcon",
		}

		test.RegisterClientNotification(notificationID)
		test.CreateNewTemplate(createdTemplate)
		test.AssignTemplateToClient(clientID)
		test.ConfirmClientTemplateAssignment(clientID)
		test.AssignTemplateToNotification(clientID, notificationID)
		test.ConfirmNotificationTemplateAssignment(clientID, notificationID)
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
	notificationsServer servers.Notifications
	clientToken         uaa.Token
	TemplateID          string
}

func (test AssignTemplate) RegisterClientNotification(notificationID string) {
	body, err := json.Marshal(map[string]interface{}{
		"source_name": "Notifications Sender",
		"notifications": map[string]interface{}{
			notificationID: map[string]interface{}{
				"description": "Acceptance Test",
				"critical":    true,
			},
		},
	})
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest("PUT", test.notificationsServer.NotificationsPath(), bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+test.clientToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	// Confirm response status code looks ok
	Expect(response.StatusCode).To(Equal(http.StatusNoContent))
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

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var JSON struct {
		TemplateID string `json:"template_id"`
	}

	err = json.Unmarshal(body, &JSON)
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
