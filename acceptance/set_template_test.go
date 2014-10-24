package acceptance

import (
    "bytes"
    "encoding/json"
    "io/ioutil"
    "net/http"

    "github.com/cloudfoundry-incubator/notifications/acceptance/servers"
    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Templates PUT Endpoint", func() {
    BeforeEach(func() {
        TruncateTables()
    })

    It("allows a user to set body templates", func() {
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
        clientID := "notifications-sender"
        env := config.NewEnvironment()
        uaaClient := uaa.NewUAA("", env.UAAHost, clientID, "secret", "")
        clientToken, err := uaaClient.GetClientToken()
        if err != nil {
            panic(err)
        }

        test := SetTemplates{}
        test.SetDefaultTemplate(notificationsServer, clientToken)
        test.GetOverriddenDefaultTemplate(notificationsServer, clientToken)
    })
})

type SetTemplates struct{}

func (t SetTemplates) SetDefaultTemplate(notificationsServer servers.Notifications, clientToken uaa.Token) {
    jsonBody := []byte(`{"text":"rulebook", "html":"<p>follow the rules</p>"}`)
    request, err := http.NewRequest("PUT", notificationsServer.SpaceTemplatePath(), bytes.NewBuffer(jsonBody))
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

func (t SetTemplates) GetOverriddenDefaultTemplate(notificationsServer servers.Notifications, clientToken uaa.Token) {
    request, err := http.NewRequest("GET", notificationsServer.SpaceTemplatePath(), bytes.NewBuffer([]byte{}))
    if err != nil {
        panic(err)
    }

    request.Header.Set("Authorization", "Bearer "+clientToken.Access)

    response, err := http.DefaultClient.Do(request)
    if err != nil {
        panic(err)
    }

    // Confirm response status code is a 204
    Expect(response.StatusCode).To(Equal(http.StatusOK))

    // Confirm we got the correct template info
    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        panic(err)
    }

    responseJSON := models.Template{}
    err = json.Unmarshal(body, &responseJSON)
    if err != nil {
        panic(err)
    }

    Expect(responseJSON.Text).To(Equal("rulebook"))
    Expect(responseJSON.HTML).To(Equal("<p>follow the rules</p>"))
    Expect(responseJSON.Overridden).To(BeTrue())
}
