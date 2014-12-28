package acceptance

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/acceptance/servers"
	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Default Template", func() {
	var t TemplateDefault
	var env application.Environment

	BeforeEach(func() {
		clientID := "notifications-admin"
		env = application.NewEnvironment()
		uaaClient := uaa.NewUAA("", env.UAAHost, clientID, "secret", "")
		clientToken, err := uaaClient.GetClientToken()
		if err != nil {
			panic(err)
		}

		t = TemplateDefault{
			notificationsServer: Servers.Notifications,
			clientToken:         clientToken,
		}
	})

	It("can retrieve the default template", func() {
		defaultTemplate := models.Template{
			ID:       "default",
			Name:     "Defaultimus Prime",
			Subject:  "Robot Subject: {{.Subject}}",
			HTML:     "<p>Default html</p>",
			Text:     "default text",
			Metadata: "{}",
		}

		conn := models.NewDatabase(env.DatabaseURL, env.ModelMigrationsDir).Connection()
		err := conn.Insert(&defaultTemplate)
		if err != nil {
			panic(err)
		}

		t.GetDefaultTemplate()
	})

	It("can edit the default template", func() {
		defaultTemplate := models.Template{
			ID:       "default",
			Name:     "Defaultimus Prime",
			Subject:  "Robot Subject: {{.Subject}}",
			HTML:     "<p>Default html</p>",
			Text:     "default text",
			Metadata: "{}",
		}

		conn := models.NewDatabase(env.DatabaseURL, env.ModelMigrationsDir).Connection()
		err := conn.Insert(&defaultTemplate)
		if err != nil {
			panic(err)
		}

		t.EditDefaultTemplate()
		t.ConfirmEditedDefaultTemplate()
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

func (t TemplateDefault) EditDefaultTemplate() {
	body := `{
		"name": "A Whole New Template",
		"subject": "Updated: {{.Subject}}",
		"html": "<h1>Updated!!!</h1>",
		"text": "Updated!!!",
		"metadata": {
			"smurf":"favorite"
		}
	}`

	request, err := http.NewRequest("PUT", t.notificationsServer.DefaultTemplatePath(), strings.NewReader(body))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+t.clientToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	Expect(response.StatusCode).To(Equal(http.StatusNoContent))
}

func (t TemplateDefault) ConfirmEditedDefaultTemplate() {
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
		"name": "A Whole New Template",
		"subject": "Updated: {{.Subject}}",
		"html": "<h1>Updated!!!</h1>",
		"text": "Updated!!!",
		"metadata": {
			"smurf":"favorite"
		}
	}`))

}
