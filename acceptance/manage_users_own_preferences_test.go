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
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/web/services"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Preferences Endpoint", func() {
	It("user unsubscribes from a notification", func() {
		userGUID := "user-123"
		clientToken := GetClientTokenFor("notifications-sender")
		userToken := GetUserTokenFor("user-123-code")

		test := ManageUsersOwnPreferences{
			client:              support.NewClient(Servers.Notifications),
			notificationsServer: Servers.Notifications,
			smtpServer:          Servers.SMTP,
			clientToken:         clientToken,
			userToken:           userToken,
			userGUID:            userGUID,
		}

		test.RegisterClientNotifications()
		test.SendNotificationToUser()
		test.RetrieveUserPreferences()

		// Notification Unsubscribe
		test.UnsubscribeFromNotification()
		test.ConfirmUserUnsubscribed()
		test.ConfirmsUnsubscribedNotificationsAreNotReceived()
		test.ResubscribeToNotification()
		test.ConfirmUserResubscribed()

		// Global Unsubscribe
		test.GlobalUnsubscribe()
		test.ConfirmGlobalUnsubscribe()
		test.ConfirmUserDoesNotReceiveNotificationsGlobal()
		test.UndoGlobalUnsubscribe()
		test.ReConfirmUserUnsubscribed()
		test.ConfirmUserReceivesNotificationsGlobal()
	})

})

type ManageUsersOwnPreferences struct {
	client              *support.Client
	notificationsServer servers.Notifications
	smtpServer          *servers.SMTP
	clientToken         uaa.Token
	userToken           uaa.Token
	userGUID            string
}

// Make request to /registation
func (t ManageUsersOwnPreferences) RegisterClientNotifications() {
	code, err := t.client.Notifications.Register(t.clientToken.Access, support.RegisterClient{
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
func (t ManageUsersOwnPreferences) SendNotificationToUser() {
	body, err := json.Marshal(map[string]string{
		"kind_id": "unsubscribe-acceptance-test",
		"html":    "<p>this is an acceptance test</p>",
		"subject": "my-special-subject",
	})
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest("POST", t.notificationsServer.UsersPath(t.userGUID), bytes.NewBuffer(body))
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
	responseItem := responseJSON[0]
	Expect(responseItem["status"]).To(Equal("queued"))
	Expect(responseItem["recipient"]).To(Equal(t.userGUID))
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
	Expect(data).To(ContainElement("X-CF-Notification-ID: " + responseItem["notification_id"]))
	Expect(data).To(ContainElement("Subject: CF Notification: my-special-subject"))
	Expect(data).To(ContainElement("        <p>this is an acceptance test</p>"))
}

// Make a GET request to /user_preferences
func (t ManageUsersOwnPreferences) RetrieveUserPreferences() {
	request, err := http.NewRequest("GET", t.notificationsServer.UserPreferencesPath(), nil)
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+t.userToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	// Confirm the request response looks correct
	Expect(response.StatusCode).To(Equal(http.StatusOK))

	prefsResponseJSON := services.PreferencesBuilder{}
	err = json.NewDecoder(response.Body).Decode(&prefsResponseJSON)
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
func (t ManageUsersOwnPreferences) UnsubscribeFromNotification() {
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

	request, err := http.NewRequest("PATCH", t.notificationsServer.UserPreferencesPath(), bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+t.userToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	Expect(response.StatusCode).To(Equal(http.StatusNoContent))
}

// Make a GET request to /user_preferences
func (t ManageUsersOwnPreferences) ConfirmUserUnsubscribed() {
	request, err := http.NewRequest("GET", t.notificationsServer.UserPreferencesPath(), nil)
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+t.userToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}
	// Confirm the request response looks correct
	Expect(response.StatusCode).To(Equal(http.StatusOK))

	prefsResponseJSON := services.PreferencesBuilder{}
	err = json.NewDecoder(response.Body).Decode(&prefsResponseJSON)
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
func (t ManageUsersOwnPreferences) ConfirmsUnsubscribedNotificationsAreNotReceived() {
	t.smtpServer.Reset() //clears deliveries

	body, err := json.Marshal(map[string]string{
		"kind_id": "unsubscribe-acceptance-test",
		"html":    "<p>this is an acceptance test</p>",
		"subject": "my-special-subject",
	})
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest("POST", t.notificationsServer.UsersPath(t.userGUID), bytes.NewBuffer(body))
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
	responseItem := responseJSON[0]
	Expect(responseItem["status"]).To(Equal("queued"))
	Expect(responseItem["recipient"]).To(Equal(t.userGUID))
	Expect(GUIDRegex.MatchString(responseItem["notification_id"])).To(BeTrue())

	// Confirm the email message was never delivered
	Consistently(func() int {
		return len(t.smtpServer.Deliveries)
	}, 1*time.Second).Should(Equal(0))
}

// Make PATCH request to /user_preferences
func (t ManageUsersOwnPreferences) ResubscribeToNotification() {
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

	request, err := http.NewRequest("PATCH", t.notificationsServer.UserPreferencesPath(), bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+t.userToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	Expect(response.StatusCode).To(Equal(http.StatusNoContent))
}

// Make a GET request to /user_preferences
func (t ManageUsersOwnPreferences) ConfirmUserResubscribed() {
	request, err := http.NewRequest("GET", t.notificationsServer.UserPreferencesPath(), nil)
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+t.userToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	// Confirm the request response looks correct
	Expect(response.StatusCode).To(Equal(http.StatusOK))

	prefsResponseJSON := services.PreferencesBuilder{}
	err = json.NewDecoder(response.Body).Decode(&prefsResponseJSON)
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

func (t ManageUsersOwnPreferences) GlobalUnsubscribe() {
	requestBodyPayload := map[string]interface{}{
		"global_unsubscribe": true,
		"clients":            map[string]interface{}{},
	}

	body, err := json.Marshal(requestBodyPayload)
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest("PATCH", t.notificationsServer.UserPreferencesPath(), bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+t.userToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	// Confirm the request response looks correct
	Expect(response.StatusCode).To(Equal(http.StatusNoContent))
}

func (t ManageUsersOwnPreferences) ConfirmGlobalUnsubscribe() {
	request, err := http.NewRequest("GET", t.notificationsServer.UserPreferencesPath(), nil)
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+t.userToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	// Confirm the request response looks correct
	Expect(response.StatusCode).To(Equal(http.StatusOK))

	prefsResponseJSON := services.PreferencesBuilder{}
	err = json.NewDecoder(response.Body).Decode(&prefsResponseJSON)
	if err != nil {
		panic(err)
	}

	Expect(prefsResponseJSON.GlobalUnsubscribe).To(BeTrue())
}

func (t ManageUsersOwnPreferences) ConfirmUserDoesNotReceiveNotificationsGlobal() {
	t.smtpServer.Reset()

	body, err := json.Marshal(map[string]string{
		"kind_id": "acceptance-test",
		"html":    "<p>this is an acceptance test</p>",
		"subject": "my-special-subject",
	})
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest("POST", t.notificationsServer.UsersPath(t.userGUID), bytes.NewBuffer(body))
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
	responseItem := responseJSON[0]
	Expect(responseItem["status"]).To(Equal("queued"))
	Expect(responseItem["recipient"]).To(Equal(t.userGUID))
	Expect(GUIDRegex.MatchString(responseItem["notification_id"])).To(BeTrue())

	// Confirm the email message never gets delivered
	Consistently(func() int {
		return len(t.smtpServer.Deliveries)
	}, 1*time.Second).Should(Equal(0))
}

func (t ManageUsersOwnPreferences) UndoGlobalUnsubscribe() {
	requestBodyPayload := map[string]interface{}{
		"global_unsubscribe": false,
		"clients":            map[string]interface{}{},
	}

	body, err := json.Marshal(requestBodyPayload)
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest("PATCH", t.notificationsServer.UserPreferencesPath(), bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+t.userToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	// Confirm the request response looks correct
	Expect(response.StatusCode).To(Equal(http.StatusNoContent))
}

func (t ManageUsersOwnPreferences) ReConfirmUserUnsubscribed() {
	request, err := http.NewRequest("GET", t.notificationsServer.UserPreferencesPath(), nil)
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+t.userToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	// Confirm the request response looks correct
	Expect(response.StatusCode).To(Equal(http.StatusOK))

	prefsResponseJSON := services.PreferencesBuilder{}
	err = json.NewDecoder(response.Body).Decode(&prefsResponseJSON)
	if err != nil {
		panic(err)
	}

	Expect(prefsResponseJSON.GlobalUnsubscribe).To(BeFalse())
}

func (t ManageUsersOwnPreferences) ConfirmUserReceivesNotificationsGlobal() {
	t.smtpServer.Reset()

	body, err := json.Marshal(map[string]string{
		"kind_id": "acceptance-test",
		"html":    "<p>this is an acceptance test</p>",
		"subject": "my-special-subject",
	})
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest("POST", t.notificationsServer.UsersPath(t.userGUID), bytes.NewBuffer(body))
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
	responseItem := responseJSON[0]
	Expect(responseItem["status"]).To(Equal("queued"))
	Expect(responseItem["recipient"]).To(Equal(t.userGUID))
	Expect(GUIDRegex.MatchString(responseItem["notification_id"])).To(BeTrue())

	// Confirm the email message gets delivered
	Eventually(func() int {
		return len(t.smtpServer.Deliveries)
	}, 1*time.Second).Should(Equal(1))
}
