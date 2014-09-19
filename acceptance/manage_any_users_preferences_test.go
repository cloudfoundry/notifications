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

    It("client unsubscribes a user from a notification", func() {
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
        env := config.NewEnvironment()
        uaaClient := uaa.NewUAA("", env.UAAHost, "notifications-sender", "secret", "")
        clientToken, err := uaaClient.GetClientToken()
        if err != nil {
            panic(err)
        }

        userGUID := "user-123"
        // Retrieve User UAA token
        userToken, err := uaaClient.Exchange(userGUID + "-code")
        if err != nil {
            panic(err)
        }

        //GRAB USER GUID FROM USER TOKEN

        test := ManageArbitraryUsersPreferences{}

        test.RegisterClientNotifications(notificationsServer, clientToken)
        test.SendNotificationToUser(notificationsServer, clientToken, smtpServer)
        test.RetrieveUserPreferences(notificationsServer, userToken)
        test.UnsubscribeFromNotification(notificationsServer, clientToken, userGUID)
        test.ConfirmUserUnsubscribed(notificationsServer, userToken)
    })
})

type ManageArbitraryUsersPreferences struct{}

// Make request to /registation
func (t ManageArbitraryUsersPreferences) RegisterClientNotifications(notificationsServer servers.Notifications, clientToken uaa.Token) {
    body, err := json.Marshal(map[string]interface{}{
        "source_description": "Notifications Sender",
        "kinds": []map[string]string{
            {
                "id":          "acceptance-test",
                "description": "Acceptance Test",
            },
            {
                "id":          "unsubscribe-acceptance-test",
                "description": "Unsubscribe Acceptance Test",
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

// Make request to /users/:guid
func (t ManageArbitraryUsersPreferences) SendNotificationToUser(notificationsServer servers.Notifications, clientToken uaa.Token, smtpServer *servers.SMTP) {
    body, err := json.Marshal(map[string]string{
        "kind_id": "unsubscribe-acceptance-test",
        "html":    "<p>this is an acceptance test</p>",
        "subject": "my-special-subject",
    })
    if err != nil {
        panic(err)
    }

    request, err := http.NewRequest("POST", notificationsServer.UsersPath("user-123"), bytes.NewBuffer(body))
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
    Expect(responseItem["recipient"]).To(Equal("user-123"))
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
    Expect(data).To(ContainElement("X-CF-Notification-ID: " + responseItem["notification_id"]))
    Expect(data).To(ContainElement("Subject: CF Notification: my-special-subject"))
    Expect(data).To(ContainElement(`        <p>The following "Unsubscribe Acceptance Test" notification was sent to you directly by the "Notifications Sender" component of Cloud Foundry:</p>`))
    Expect(data).To(ContainElement("<p>this is an acceptance test</p>"))
}

func (t ManageArbitraryUsersPreferences) RetrieveUserPreferences(notificationsServer servers.Notifications, userToken uaa.Token) {
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

    node := prefsResponseJSON["notifications-sender"]["acceptance-test"]
    Expect(node).To(Equal(map[string]interface{}{
        "email":              true,
        "kind_description":   "Acceptance Test",
        "source_description": "Notifications Sender",
    }))

    node = prefsResponseJSON["notifications-sender"]["unsubscribe-acceptance-test"]
    Expect(node).To(Equal(map[string]interface{}{
        "email":              true,
        "kind_description":   "Unsubscribe Acceptance Test",
        "source_description": "Notifications Sender",
    }))
}

func (t ManageArbitraryUsersPreferences) UnsubscribeFromNotification(notificationsServer servers.Notifications, clientToken uaa.Token, userGUID string) {
    builder := services.NewPreferencesBuilder()
    builder.Add(models.Preference{
        ClientID: "notifications-sender",
        KindID:   "unsubscribe-acceptance-test",
        Email:    false,
    })

    body, err := json.Marshal(builder)
    if err != nil {
        panic(err)
    }

    request, err := http.NewRequest("PATCH", notificationsServer.SpecificUserPreferencesPath(userGUID), bytes.NewBuffer(body))
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

}

func (t ManageArbitraryUsersPreferences) ConfirmUserUnsubscribed(notificationsServer servers.Notifications, userToken uaa.Token) {
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

    node := prefsResponseJSON["notifications-sender"]["acceptance-test"]
    Expect(node).To(Equal(map[string]interface{}{
        "email":              true,
        "kind_description":   "Acceptance Test",
        "source_description": "Notifications Sender",
    }))

    node = prefsResponseJSON["notifications-sender"]["unsubscribe-acceptance-test"]
    Expect(node).To(Equal(map[string]interface{}{
        "email":              false,
        "kind_description":   "Unsubscribe Acceptance Test",
        "source_description": "Notifications Sender",
    }))
}
