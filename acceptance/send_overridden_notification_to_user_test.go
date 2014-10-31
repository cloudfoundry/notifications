package acceptance

import (
    "bytes"
    "encoding/json"
    "fmt"
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

var _ = Describe("Send a notification to user with overridden template", func() {
    BeforeEach(func() {
        TruncateTables()
    })

    It("send a notification to user", func() {
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

        textTemplate := "text"
        htmlTemplate := "<p>html</p>"
        t := SendOverriddenNotificationToUser{}
        t.OverrideClientUserTemplate(notificationsServer, clientToken, textTemplate, htmlTemplate)
        t.SendNotificationToUser(notificationsServer, clientToken, smtpServer, textTemplate, htmlTemplate)
    })
})

type SendOverriddenNotificationToUser struct{}

func (t SendOverriddenNotificationToUser) OverrideClientUserTemplate(notificationsServer servers.Notifications, clientToken uaa.Token, textTemplate, htmlTemplate string) {
    jsonBody := []byte(fmt.Sprintf(`{"text":"%s", "html":"%s"}`, textTemplate, htmlTemplate))
    request, err := http.NewRequest("PUT", notificationsServer.TemplatePath("notifications-sender.user_body"), bytes.NewBuffer(jsonBody))
    if err != nil {
        panic(err)
    }

    request.Header.Set("Authorization", "Bearer "+clientToken.Access)

    response, err := http.DefaultClient.Do(request)
    if err != nil {
        panic(err)
    }

    // Confirm response status code is a 204
    Expect(response.StatusCode).To(Equal(http.StatusNoContent))
}

func (t SendOverriddenNotificationToUser) SendNotificationToUser(notificationsServer servers.Notifications, clientToken uaa.Token,
    smtpServer *servers.SMTP, text, html string) {

    body, err := json.Marshal(map[string]string{
        "kind_id": "acceptance-test",
        "html":    "<p>this is an acceptance%40test</p>",
        "text":    "the acceptance text",
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

    Expect(data).To(ContainElement(text))
    Expect(data).To(ContainElement("        <p>html</p>"))
}
