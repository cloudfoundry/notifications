package acceptance

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/acceptance/servers"
	"github.com/cloudfoundry-incubator/notifications/application"
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
		env := application.NewEnvironment()
		uaaClient := uaa.NewUAA("", env.UAAHost, clientID, "secret", "")
		clientToken, err := uaaClient.GetClientToken()
		if err != nil {
			panic(err)
		}

		t := SetTemplates{
			notificationsServer: notificationsServer,
			clientToken:         clientToken,
			text:                "rulebook",
			html:                "<p>follow the rules</p>",
		}
		t.SetDefaultSpaceTemplate()
	})
})

type SetTemplates struct {
	notificationsServer servers.Notifications
	clientToken         uaa.Token
	text                string
	html                string
}

func (t SetTemplates) SetDefaultSpaceTemplate() {
	jsonBody := []byte(fmt.Sprintf(`{"text":"%s", "html":"%s"}`, t.text, t.html))
	request, err := http.NewRequest("PUT", t.notificationsServer.DeprecatedTemplatePath(models.SpaceBodyTemplateName), bytes.NewBuffer(jsonBody))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+t.clientToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	// Confirm response status code is a 204
	Expect(response.StatusCode).To(Equal(http.StatusNoContent))
}
