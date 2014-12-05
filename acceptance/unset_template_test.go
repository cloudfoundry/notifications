package acceptance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/acceptance/servers"
	"github.com/cloudfoundry-incubator/notifications/config"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = XDescribe("Templates DELETE Endpoint", func() {
	BeforeEach(func() {
		TruncateTables()
	})

	It("allows a user to unset body templates", func() {
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

		text := "rulebook, reading"
		html := "<p>follow the rules</p>"
		test := DeleteTemplates{}
		test.SetTemplates(notificationsServer, clientToken, text, html)
		test.DeleteTemplates(notificationsServer, clientToken)
		test.GetTemplates(notificationsServer, clientToken)
	})
})

type DeleteTemplates struct{}

func (t DeleteTemplates) SetTemplates(notificationsServer servers.Notifications, clientToken uaa.Token, text, html string) {
	jsonBody := []byte(fmt.Sprintf(`{"text":"%s", "html":"%s"}`, text, html))
	request, err := http.NewRequest("PUT", notificationsServer.TemplatePath(models.SpaceBodyTemplateName), bytes.NewBuffer(jsonBody))
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

func (t DeleteTemplates) DeleteTemplates(notificationsServer servers.Notifications, clientToken uaa.Token) {
	request, err := http.NewRequest("DELETE", notificationsServer.TemplatePath(models.SpaceBodyTemplateName), bytes.NewBuffer([]byte(``)))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+clientToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	Expect(response.StatusCode).To(Equal(http.StatusNoContent))
}

func (t DeleteTemplates) GetTemplates(notificationsServer servers.Notifications, clientToken uaa.Token) {
	request, err := http.NewRequest("GET", notificationsServer.TemplatePath(models.SpaceBodyTemplateName), bytes.NewBuffer([]byte(``)))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+clientToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

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

	Expect(responseJSON.Text).To(Equal(`Hello {{.To}},

The following "{{.KindDescription}}" notification was sent to you by the "{{.SourceDescription}}"
component of Cloud Foundry because you are a member of the "{{.Space}}" space
in the "{{.Organization}}" organization:

{{.Text}}

This message was sent from {{.From}} and can be replied to at {{.ReplyTo}}. The
notification can be identified with the {{.MessageID}} identifier and was sent
with the {{.ClientID}} UAA client. The notification can be unsubscribed from
using the "{{.UnsubscribeID}}" unsubscribe token.
`))

	Expect(responseJSON.HTML).To(Equal(`<p>Hello {{.To}},</p>

<p>The following "{{.KindDescription}}" notification was sent to you by the "{{.SourceDescription}}"
    component of Cloud Foundry because you are a member of the "{{.Space}}" space
    in the "{{.Organization}}" organization:</p>

{{.HTML}}

<p>This message was sent from {{.From}} and can be replied to at {{.ReplyTo}}. The
    notification can be identified with the {{.MessageID}} identifier and was sent
    with the {{.ClientID}} UAA client. The notification can be unsubscribed from
    using the "{{.UnsubscribeID}}" unsubscribe token.</p>
`))

	Expect(responseJSON.Overridden).To(BeFalse())
}
