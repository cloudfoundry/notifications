package acceptance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/cloudfoundry-incubator/notifications/acceptance/servers"
	"github.com/cloudfoundry-incubator/notifications/config"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create a new template", func() {
	BeforeEach(func() {
		TruncateTables()

		env := config.NewEnvironment()
		migrationsPath := path.Join(env.RootPath, env.ModelMigrationsDir)
		models.NewDatabase(env.DatabaseURL, migrationsPath) // this is the "database" variable
	})

	It("allows a user to create a new template", func() {
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

		test := TemplatesCRUD{}
		name := "Star Wars"
		subject := "Awesomeness"
		html := "<p>Millenium Falcon</p>"
		text := "Millenium Falcon"

		test.CreateNewTemplate(notificationsServer, clientToken, name, subject, html, text)
		test.GetTemplate(notificationsServer, clientToken, name, subject, html, text)
	})
})

type TemplatesCRUD struct{}

func (test TemplatesCRUD) CreateNewTemplate(notificationsServer servers.Notifications, clientToken uaa.Token, name, subject, html, text string) {
	jsonBody := []byte(fmt.Sprintf(`{"name":"%s", "subject":"%s", "html":"%s", "text":"%s"}`, name, subject, html, text))
	request, err := http.NewRequest("POST", notificationsServer.TemplatesBasePath(), bytes.NewBuffer(jsonBody))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+clientToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var responseMap map[string]string
	err = json.Unmarshal(body, &responseMap)
	if err != nil {
		panic(err)
	}

	Expect(response.StatusCode).To(Equal(http.StatusCreated))
	Expect(responseMap).To(HaveKey("template-id"))
	Expect(responseMap["template-id"]).ToNot(BeNil())
}

func (test TemplatesCRUD) GetTemplate(notificationsServer servers.Notifications, clientToken uaa.Token, name, subject, html, text string) {
	request, err := http.NewRequest("GET", notificationsServer.TemplatePath("guid"), bytes.NewBuffer([]byte{}))
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

	Expect(responseJSON.Name).To(Equal(name))
	Expect(responseJSON.Subject).To(Equal(subject))
	Expect(responseJSON.Text).To(Equal(text))
	Expect(responseJSON.HTML).To(Equal(html))
}
