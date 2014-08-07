package acceptance

import (
    "bytes"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/acceptance/servers"
    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Sending notifications to all users in a space", func() {
    It("sends a notification to each user in a space", func() {
        // Boot Fake SMTP Server
        smtpServer := servers.NewSMTPServer()
        smtpServer.Boot()

        // Boot Fake UAA Server
        uaaServer := servers.NewUAAServer()
        uaaServer.Boot()
        defer uaaServer.Close()

        // Boot Fake CC Server
        ccServer := servers.NewCCServer()
        ccServer.Boot()
        defer ccServer.Close()

        // Boot Real Notifications Server
        notificationsServer := servers.NewNotificationsServer()
        notificationsServer.Boot()
        defer notificationsServer.Close()

        // Make request to /users/:guid
        body, err := json.Marshal(map[string]string{
            "kind":    "space-test",
            "text":    "this is a space test",
            "subject": "space-subject",
        })
        request, err := http.NewRequest("POST", notificationsServer.SpacesPath("space-123"), bytes.NewBuffer(body))
        if err != nil {
            panic(err)
        }

        env := config.NewEnvironment()
        uaaClient := uaa.NewUAA("", env.UAAHost, "notifications-sender", "secret", "")
        token, err := uaaClient.GetClientToken()
        if err != nil {
            panic(err)
        }
        request.Header.Set("Authorization", "Bearer "+token.Access)

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

        Expect(len(responseJSON)).To(Equal(3))

        indexedResponses := map[string]map[string]string{}
        for _, resp := range responseJSON {
            indexedResponses[resp["recipient"]] = resp
        }

        responseItem := indexedResponses["user-456"]
        Expect(responseItem["recipient"]).To(Equal("user-456"))
        Expect(responseItem["status"]).To(Equal("delivered"))
        Expect(GUIDRegex.MatchString(responseItem["notification_id"])).To(BeTrue())

        responseItem = indexedResponses["user-789"]
        Expect(responseItem["recipient"]).To(Equal("user-789"))
        Expect(responseItem["status"]).To(Equal("noaddress"))
        Expect(responseItem["notification_id"]).To(Equal(""))

        responseItem = indexedResponses["user-000"]
        Expect(responseItem["recipient"]).To(Equal("user-000"))
        Expect(responseItem["status"]).To(Equal("notfound"))
        Expect(responseItem["notification_id"]).To(Equal(""))

        // Confirm the email message was delivered correctly
        Expect(len(smtpServer.Deliveries)).To(Equal(1))
        delivery := smtpServer.Deliveries[0]

        Expect(delivery.Sender).To(Equal(env.Sender))
        Expect(delivery.Recipients).To(Equal([]string{"user-456@example.com"}))

        data := strings.Split(string(delivery.Data), "\n")
        Expect(data).To(ContainElement("X-CF-Client-ID: notifications-sender"))
        Expect(data).To(ContainElement("X-CF-Notification-ID: " + indexedResponses["user-456"]["notification_id"]))
        Expect(data).To(ContainElement("Subject: CF Notification: space-subject"))
        Expect(data).To(ContainElement("this is a space test"))
    })
})
