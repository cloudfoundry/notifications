package acceptance

import (
	"bytes"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/acceptance/servers"
	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Default Template", func() {
	BeforeEach(func() {
		TruncateTables()
	})

	It("retrieves the default template from the database", func() {
		smtpServer := servers.NewSMTP()
		smtpServer.Boot()

		uaaServer := servers.NewUAA()
		uaaServer.Boot()
		defer uaaServer.Close()

		notificationsServer := servers.NewNotifications()
		notificationsServer.Boot()
		defer notificationsServer.Close()

		clientID := "notifications-admin"
		env := application.NewEnvironment()
		uaaClient := uaa.NewUAA("", env.UAAHost, clientID, "secret", "")
		clientToken, err := uaaClient.GetClientToken()
		if err != nil {
			panic(err)
		}

		defaultTemplate := models.Template{
			ID:       "default",
			Name:     "Defaultimus Prime",
			Subject:  "Robot Subject: {{.Subject}}",
			HTML:     "<p>Default html</p>",
			Text:     "default text",
			Metadata: "{}",
		}

		conn := models.NewDatabase(env.DatabaseURL, env.ModelMigrationsDir).Connection()
		err = conn.Insert(&defaultTemplate)
		if err != nil {
			panic(err)
		}

		t := TemplateDefault{
			notificationsServer: notificationsServer,
			clientToken:         clientToken,
		}

		t.GetDefaultTemplate()
	})

})

type TemplateDefault struct {
	notificationsServer servers.Notifications
	clientToken         uaa.Token
}

func (t TemplateDefault) GetDefaultTemplate() {
	request, err := http.NewRequest("GET", t.notificationsServer.DefaultTemplatePath(), nil)
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+t.clientToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	Expect(response.StatusCode).To(Equal(http.StatusOK))

	buffer := bytes.NewBuffer([]byte{})
	_, err = buffer.ReadFrom(response.Body)
	if err != nil {
		panic(err)
	}

	Expect(buffer).To(MatchJSON(`{
		"name":     "Defaultimus Prime",
		"subject":  "Robot Subject: {{.Subject}}",
		"html":     "<p>Default html</p>",
		"text":     "default text",
		"metadata": {}
	}`))
}
