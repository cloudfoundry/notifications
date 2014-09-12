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

var _ = Describe("Send a notification to an email", func() {
    BeforeEach(func() {
        TruncateTables()
    })

    It("sends a single notification to an email", func() {
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

        // Retrieve UAA token
        env := config.NewEnvironment()
        uaaClient := uaa.NewUAA("", env.UAAHost, "notifications-sender", "secret", "")
        clientToken, err := uaaClient.GetClientToken()
        if err != nil {
            panic(err)
        }

        t := SendNotificationToEmail{}
        t.SendNotificationToEmail(notificationsServer, clientToken, smtpServer)
    })

})

type SendNotificationToEmail struct{}

func (t SendNotificationToEmail) SendNotificationToEmail(notificationsServer servers.Notifications, clientToken uaa.Token, smtpServer *servers.SMTP) {

    body, err := json.Marshal(map[string]string{
        "kind_id": "acceptance-test",
        "html":    "<p>this is an acceptance test</p>",
        "subject": "my-special-subject",
        "to":      "John User <user@example.com>",
    })

    if err != nil {
        panic(err)
    }

    request, err := http.NewRequest("POST", notificationsServer.EmailPath(), bytes.NewBuffer(body))
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
    Expect(responseItem["email"]).To(Equal("user@example.com"))
    Expect(GUIDRegex.MatchString(responseItem["notification_id"])).To(BeTrue())

    // Confirm the email message was delivered correctly
    Eventually(func() int {
        return len(smtpServer.Deliveries)
    }, 5*time.Second).Should(Equal(1))
    delivery := smtpServer.Deliveries[0]

    env := config.NewEnvironment()
    Expect(delivery.Sender).To(Equal(env.Sender))
    Expect(delivery.Recipients).To(Equal([]string{"user@example.com"}))

    data := strings.Split(string(delivery.Data), "\n")
    Expect(data).To(ContainElement("X-CF-Client-ID: notifications-sender"))
    Expect(data).To(ContainElement("X-CF-Notification-ID: " + responseItem["notification_id"]))
    Expect(data).To(ContainElement("Subject: CF Notification: my-special-subject"))
    Expect(data).To(ContainElement("        the template"))
    Expect(data).To(ContainElement("<p>this is an acceptance test</p>"))
}
