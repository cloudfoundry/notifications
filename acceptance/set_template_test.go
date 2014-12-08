package acceptance

import (
	"bytes"
	"fmt"
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
		text := "rulebook"
		html := "<p>follow the rules</p>"
		test.SetDefaultSpaceTemplate(notificationsServer, clientToken, text, html)
	})
})

type SetTemplates struct{}

func (t SetTemplates) SetDefaultSpaceTemplate(notificationsServer servers.Notifications, clientToken uaa.Token, text, html string) {
	jsonBody := []byte(fmt.Sprintf(`{"text":"%s", "html":"%s"}`, text, html))
	request, err := http.NewRequest("PUT", notificationsServer.DeprecatedTemplatePath(models.SpaceBodyTemplateName), bytes.NewBuffer(jsonBody))
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
