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
	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/web/services"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Preferences Endpoint", func() {
	BeforeEach(func() {
		TruncateTables()
	})

	It("user unsubscribes from a notification", func() {
		userGUID := "user-123"

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
		env := application.NewEnvironment()
		uaaClient := uaa.NewUAA("", env.UAAHost, "notifications-sender", "secret", "")
		clientToken, err := uaaClient.GetClientToken()
		if err != nil {
			panic(err)
		}

		// Retrieve User UAA token
		userToken, err := uaaClient.Exchange("user-123-code")
		if err != nil {
			panic(err)
		}

		test := ManageUsersOwnPreferences{
			client: support.NewClient(notificationsServer),
		}

		test.RegisterClientNotifications(notificationsServer, clientToken)
		test.SendNotificationToUser(notificationsServer, clientToken, userGUID, smtpServer)
		test.RetrieveUserPreferences(notificationsServer, userToken)

		// Notification Unsubscribe
		test.UnsubscribeFromNotification(notificationsServer, userToken)
		test.ConfirmUserUnsubscribed(notificationsServer, userToken)
		test.ConfirmsUnsubscribedNotificationsAreNotReceived(notificationsServer, clientToken, userGUID, smtpServer)
		test.ResubscribeToNotification(notificationsServer, userToken)
		test.ConfirmUserResubscribed(notificationsServer, userToken)

		// Global Unsubscribe
		test.GlobalUnsubscribe(notificationsServer, userToken)
		test.ConfirmGlobalUnsubscribe(notificationsServer, userToken)
		test.ConfirmUserDoesNotReceiveNotificationsGlobal(notificationsServer, clientToken, userGUID, smtpServer)
		test.UndoGlobalUnsubscribe(notificationsServer, userToken)
		test.ReConfirmUserUnsubscribed(notificationsServer, userToken)
		test.ConfirmUserReceivesNotificationsGlobal(notificationsServer, clientToken, userGUID, smtpServer)
	})

})

type ManageUsersOwnPreferences struct {
	client *support.Client
}

// Make request to /registation
func (t ManageUsersOwnPreferences) RegisterClientNotifications(notificationsServer servers.Notifications, clientToken uaa.Token) {
	code, err := t.client.Notifications.Register(clientToken.Access, support.RegisterClient{
		SourceName: "Notifications Sender",
		Notifications: map[string]support.RegisterNotification{
			"acceptance-test": {
				Description: "Acceptance Test",
			},
			"unsubscribe-acceptance-test": {
				Description: "Unsubscribe Acceptance Test",
			},
		},
	})

	Expect(err).NotTo(HaveOccurred())
	Expect(code).To(Equal(http.StatusNoContent))
}

// Make request to /users/:guid
func (t ManageUsersOwnPreferences) SendNotificationToUser(notificationsServer servers.Notifications, clientToken uaa.Token, userGUID string, smtpServer *servers.SMTP) {
	body, err := json.Marshal(map[string]string{
		"kind_id": "unsubscribe-acceptance-test",
		"html":    "<p>this is an acceptance test</p>",
		"subject": "my-special-subject",
	})
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest("POST", notificationsServer.UsersPath(userGUID), bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+clientToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	// Confirm the request response looks correct
	Expect(response.StatusCode).To(Equal(http.StatusOK))

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	responseJSON := []map[string]string{}
	err = json.Unmarshal(body, &responseJSON)
	if err != nil {
		panic(err)
	}

	Expect(len(responseJSON)).To(Equal(1))
	responseItem := responseJSON[0]
	Expect(responseItem["status"]).To(Equal("queued"))
	Expect(responseItem["recipient"]).To(Equal(userGUID))
	Expect(GUIDRegex.MatchString(responseItem["notification_id"])).To(BeTrue())

	// Confirm the email message was delivered correctly
	Eventually(func() int {
		return len(smtpServer.Deliveries)
	}, 5*time.Second).Should(Equal(1))
	delivery := smtpServer.Deliveries[0]

	env := application.NewEnvironment()
	Expect(delivery.Sender).To(Equal(env.Sender))
	Expect(delivery.Recipients).To(Equal([]string{"user-123@example.com"}))

	data := strings.Split(string(delivery.Data), "\n")
	Expect(data).To(ContainElement("X-CF-Client-ID: notifications-sender"))
	Expect(data).To(ContainElement("X-CF-Notification-ID: " + responseItem["notification_id"]))
	Expect(data).To(ContainElement("Subject: CF Notification: my-special-subject"))
	Expect(data).To(ContainElement(`<p>The following "Unsubscribe Acceptance Test" notification was sent to you directly by the`))
	Expect(data).To(ContainElement(`    "Notifications Sender" component of Cloud Foundry:</p>`))
	Expect(data).To(ContainElement("<p>this is an acceptance test</p>"))
}

// Make a GET request to /user_preferences
func (t ManageUsersOwnPreferences) RetrieveUserPreferences(notificationsServer servers.Notifications, userToken uaa.Token) {
	request, err := http.NewRequest("GET", notificationsServer.UserPreferencesPath(), nil)
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+userToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	// Confirm the request response looks correct
	Expect(response.StatusCode).To(Equal(http.StatusOK))

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	prefsResponseJSON := services.PreferencesBuilder{}
	err = json.Unmarshal(body, &prefsResponseJSON)
	if err != nil {
		panic(err)
	}

	node := prefsResponseJSON.Clients["notifications-sender"]["acceptance-test"]
	Expect(node.Email).To(Equal(&TRUE))
	Expect(node.KindDescription).To(Equal("Acceptance Test"))
	Expect(node.SourceDescription).To(Equal("Notifications Sender"))
	Expect(node.Count).To(Equal(0))

	node = prefsResponseJSON.Clients["notifications-sender"]["unsubscribe-acceptance-test"]
	Expect(node.Email).To(Equal(&TRUE))
	Expect(node.KindDescription).To(Equal("Unsubscribe Acceptance Test"))
	Expect(node.SourceDescription).To(Equal("Notifications Sender"))
	Expect(node.Count).To(Equal(1))

}

// Make a PATCH request to /user_preferences
func (t ManageUsersOwnPreferences) UnsubscribeFromNotification(notificationsServer servers.Notifications, userToken uaa.Token) {
	builder := services.NewPreferencesBuilder()
	builder.Add(models.Preference{
		ClientID: "notifications-sender",
		KindID:   "unsubscribe-acceptance-test",
		Email:    false,
		Count:    23,
	})

	body, err := json.Marshal(builder)
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest("PATCH", notificationsServer.UserPreferencesPath(), bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+userToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	Expect(response.StatusCode).To(Equal(http.StatusNoContent))
}

// Make a GET request to /user_preferences
func (t ManageUsersOwnPreferences) ConfirmUserUnsubscribed(notificationsServer servers.Notifications, userToken uaa.Token) {
	request, err := http.NewRequest("GET", notificationsServer.UserPreferencesPath(), nil)
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+userToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	// Confirm the request response looks correct
	Expect(response.StatusCode).To(Equal(http.StatusOK))

	prefsResponseJSON := services.PreferencesBuilder{}
	err = json.Unmarshal(body, &prefsResponseJSON)
	if err != nil {
		panic(err)
	}

	node := prefsResponseJSON.Clients["notifications-sender"]["acceptance-test"]
	Expect(node.Email).To(Equal(&TRUE))
	Expect(node.KindDescription).To(Equal("Acceptance Test"))
	Expect(node.SourceDescription).To(Equal("Notifications Sender"))
	Expect(node.Count).To(Equal(0))

	node = prefsResponseJSON.Clients["notifications-sender"]["unsubscribe-acceptance-test"]
	Expect(node.Email).To(Equal(&FALSE))
	Expect(node.KindDescription).To(Equal("Unsubscribe Acceptance Test"))
	Expect(node.SourceDescription).To(Equal("Notifications Sender"))
	Expect(node.Count).To(Equal(1))
}

// Make request to /users/:guid
func (t ManageUsersOwnPreferences) ConfirmsUnsubscribedNotificationsAreNotReceived(notificationsServer servers.Notifications, clientToken uaa.Token, userGUID string, smtpServer *servers.SMTP) {
	smtpServer.Reset() //clears deliveries

	body, err := json.Marshal(map[string]string{
		"kind_id": "unsubscribe-acceptance-test",
		"html":    "<p>this is an acceptance test</p>",
		"subject": "my-special-subject",
	})
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest("POST", notificationsServer.UsersPath(userGUID), bytes.NewBuffer(body))
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
	responseItem := responseJSON[0]
	Expect(responseItem["status"]).To(Equal("queued"))
	Expect(responseItem["recipient"]).To(Equal(userGUID))
	Expect(GUIDRegex.MatchString(responseItem["notification_id"])).To(BeTrue())

	// Confirm the email message was never delivered
	Consistently(func() int {
		return len(smtpServer.Deliveries)
	}, 5*time.Second).Should(Equal(0))
}

// Make PATCH request to /user_preferences
func (t ManageUsersOwnPreferences) ResubscribeToNotification(notificationsServer servers.Notifications, userToken uaa.Token) {
	builder := services.NewPreferencesBuilder()
	builder.Add(models.Preference{
		ClientID: "notifications-sender",
		KindID:   "unsubscribe-acceptance-test",
		Email:    true,
		Count:    -23,
	})

	body, err := json.Marshal(builder)
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest("PATCH", notificationsServer.UserPreferencesPath(), bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+userToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	Expect(response.StatusCode).To(Equal(http.StatusNoContent))
}

// Make a GET request to /user_preferences
func (t ManageUsersOwnPreferences) ConfirmUserResubscribed(notificationsServer servers.Notifications, userToken uaa.Token) {
	request, err := http.NewRequest("GET", notificationsServer.UserPreferencesPath(), nil)
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+userToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	// Confirm the request response looks correct
	Expect(response.StatusCode).To(Equal(http.StatusOK))

	prefsResponseJSON := services.PreferencesBuilder{}
	err = json.Unmarshal(body, &prefsResponseJSON)
	if err != nil {
		panic(err)
	}

	node := prefsResponseJSON.Clients["notifications-sender"]["acceptance-test"]
	Expect(node.Email).To(Equal(&TRUE))
	Expect(node.KindDescription).To(Equal("Acceptance Test"))
	Expect(node.SourceDescription).To(Equal("Notifications Sender"))
	Expect(node.Count).To(Equal(0))

	node = prefsResponseJSON.Clients["notifications-sender"]["unsubscribe-acceptance-test"]
	Expect(node.Email).To(Equal(&TRUE))
	Expect(node.KindDescription).To(Equal("Unsubscribe Acceptance Test"))
	Expect(node.SourceDescription).To(Equal("Notifications Sender"))
	Expect(node.Count).To(Equal(2))
}

func (t ManageUsersOwnPreferences) GlobalUnsubscribe(notificationsServer servers.Notifications, userToken uaa.Token) {
	requestBodyPayload := map[string]interface{}{
		"global_unsubscribe": true,
		"clients":            map[string]interface{}{},
	}

	body, err := json.Marshal(requestBodyPayload)
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest("PATCH", notificationsServer.UserPreferencesPath(), bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+userToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	// Confirm the request response looks correct
	Expect(response.StatusCode).To(Equal(http.StatusNoContent))
}

func (t ManageUsersOwnPreferences) ConfirmGlobalUnsubscribe(notificationsServer servers.Notifications, userToken uaa.Token) {
	request, err := http.NewRequest("GET", notificationsServer.UserPreferencesPath(), nil)
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+userToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	// Confirm the request response looks correct
	Expect(response.StatusCode).To(Equal(http.StatusOK))

	prefsResponseJSON := services.PreferencesBuilder{}
	err = json.Unmarshal(body, &prefsResponseJSON)
	if err != nil {
		panic(err)
	}

	Expect(prefsResponseJSON.GlobalUnsubscribe).To(BeTrue())
}

func (t ManageUsersOwnPreferences) ConfirmUserDoesNotReceiveNotificationsGlobal(notificationsServer servers.Notifications, clientToken uaa.Token, userGUID string, smtpServer *servers.SMTP) {
	smtpServer.Reset()

	body, err := json.Marshal(map[string]string{
		"kind_id": "acceptance-test",
		"html":    "<p>this is an acceptance test</p>",
		"subject": "my-special-subject",
	})
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest("POST", notificationsServer.UsersPath(userGUID), bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+clientToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	// Confirm the request response looks correct
	Expect(response.StatusCode).To(Equal(http.StatusOK))

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	responseJSON := []map[string]string{}
	err = json.Unmarshal(body, &responseJSON)
	if err != nil {
		panic(err)
	}

	Expect(len(responseJSON)).To(Equal(1))
	responseItem := responseJSON[0]
	Expect(responseItem["status"]).To(Equal("queued"))
	Expect(responseItem["recipient"]).To(Equal(userGUID))
	Expect(GUIDRegex.MatchString(responseItem["notification_id"])).To(BeTrue())

	// Confirm the email message never gets delivered
	Consistently(func() int {
		return len(smtpServer.Deliveries)
	}, 5*time.Second).Should(Equal(0))
}

func (t ManageUsersOwnPreferences) UndoGlobalUnsubscribe(notificationsServer servers.Notifications, userToken uaa.Token) {
	requestBodyPayload := map[string]interface{}{
		"global_unsubscribe": false,
		"clients":            map[string]interface{}{},
	}

	body, err := json.Marshal(requestBodyPayload)
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest("PATCH", notificationsServer.UserPreferencesPath(), bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+userToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	// Confirm the request response looks correct
	Expect(response.StatusCode).To(Equal(http.StatusNoContent))
}

func (t ManageUsersOwnPreferences) ReConfirmUserUnsubscribed(notificationsServer servers.Notifications, userToken uaa.Token) {
	request, err := http.NewRequest("GET", notificationsServer.UserPreferencesPath(), nil)
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+userToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	// Confirm the request response looks correct
	Expect(response.StatusCode).To(Equal(http.StatusOK))

	prefsResponseJSON := services.PreferencesBuilder{}
	err = json.Unmarshal(body, &prefsResponseJSON)
	if err != nil {
		panic(err)
	}

	Expect(prefsResponseJSON.GlobalUnsubscribe).To(BeFalse())
}

func (t ManageUsersOwnPreferences) ConfirmUserReceivesNotificationsGlobal(notificationsServer servers.Notifications, clientToken uaa.Token, userGUID string, smtpServer *servers.SMTP) {
	smtpServer.Reset()

	body, err := json.Marshal(map[string]string{
		"kind_id": "acceptance-test",
		"html":    "<p>this is an acceptance test</p>",
		"subject": "my-special-subject",
	})
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest("POST", notificationsServer.UsersPath(userGUID), bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+clientToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	// Confirm the request response looks correct
	Expect(response.StatusCode).To(Equal(http.StatusOK))

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	responseJSON := []map[string]string{}
	err = json.Unmarshal(body, &responseJSON)
	if err != nil {
		panic(err)
	}

	Expect(len(responseJSON)).To(Equal(1))
	responseItem := responseJSON[0]
	Expect(responseItem["status"]).To(Equal("queued"))
	Expect(responseItem["recipient"]).To(Equal(userGUID))
	Expect(GUIDRegex.MatchString(responseItem["notification_id"])).To(BeTrue())

	// Confirm the email message gets delivered
	Eventually(func() int {
		return len(smtpServer.Deliveries)
	}, 5*time.Second).Should(Equal(1))
}
