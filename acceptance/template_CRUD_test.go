package acceptance

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/cloudfoundry-incubator/notifications/acceptance/servers"
	"github.com/cloudfoundry-incubator/notifications/config"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/web/params"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Templates CRUD", func() {
	BeforeEach(func() {
		TruncateTables()

		env := config.NewEnvironment()
		migrationsPath := path.Join(env.RootPath, env.ModelMigrationsDir)
		models.NewDatabase(env.DatabaseURL, migrationsPath) // this is the "database" variable
	})

	It("allows a user to perform CRUD actions on a template", func() {
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
		clientID := "notifications-admin"
		env := config.NewEnvironment()
		uaaClient := uaa.NewUAA("", env.UAAHost, clientID, "secret", "")
		clientToken, err := uaaClient.GetClientToken()
		if err != nil {
			panic(err)
		}

		test := TemplatesCRUD{}
		createTemplate := params.Template{
			Name:    "Star Wars",
			Subject: "Awesomeness",
			HTML:    "<p>Millenium Falcon</p>",
			Text:    "Millenium Falcon",
		}

		getTemplate := params.Template{
			Name:    "Big Hero 6",
			Subject: "Heroes",
			HTML:    "<h1>Robots!</h1>",
			Text:    "Robots!",
		}

		test.CreateNewTemplate(notificationsServer, clientToken, createTemplate)
		test.GetTemplate(notificationsServer, clientToken, getTemplate)
	})
})

type TemplatesCRUD struct{}

func (test TemplatesCRUD) CreateNewTemplate(notificationsServer servers.Notifications, clientToken uaa.Token, template params.Template) {
	templateID, status := test.createTemplate(notificationsServer, clientToken, template)
	Expect(status).To(Equal(http.StatusCreated))
	Expect(templateID).NotTo(BeNil())
}

func (test TemplatesCRUD) GetTemplate(notificationsServer servers.Notifications, clientToken uaa.Token, getTemplate params.Template) {
	templateID, _ := test.createTemplate(notificationsServer, clientToken, getTemplate)

	request, err := http.NewRequest("GET", notificationsServer.TemplatePath(templateID), bytes.NewBuffer([]byte{}))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+clientToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	Expect(response.StatusCode).To(Equal(http.StatusOK))

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	responseJSON := params.Template{}
	err = json.Unmarshal(body, &responseJSON)
	if err != nil {
		panic(err)
	}

	Expect(responseJSON).To(Equal(getTemplate))
}

func (test TemplatesCRUD) createTemplate(notificationsServer servers.Notifications, clientToken uaa.Token, getTemplate params.Template) (string, int) {
	jsonBody, err := json.Marshal(getTemplate)
	if err != nil {
		panic(err)
	}

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

	var JSON struct {
		TemplateID string `json:"template-id"`
	}

	err = json.Unmarshal(body, &JSON)
	if err != nil {
		panic(err)
	}

	return JSON.TemplateID, response.StatusCode
}
