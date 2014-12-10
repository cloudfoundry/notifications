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

		test := TemplatesCRUD{
			notificationsServer: notificationsServer,
			clientToken:         clientToken,
		}

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

		updateTemplate := params.Template{
			Name:    "Blah",
			Subject: "More Blah",
			HTML:    "<h1>This is blahblah</h1>",
			Text:    "Blah even more",
		}

		deleteTemplate := params.Template{
			Name:    "Hungry Play",
			Subject: "Dystopian",
			HTML:    "<h1>Sad</h1>",
			Text:    "Run!!",
		}

		test.CreateNewTemplate(createTemplate)
		test.GetTemplate(getTemplate)
		test.UpdateTemplate(updateTemplate)
		test.DeleteTemplate(deleteTemplate)
	})
})

type TemplatesCRUD struct {
	notificationsServer servers.Notifications
	clientToken         uaa.Token
}

func (test TemplatesCRUD) CreateNewTemplate(template params.Template) {
	templateID, status := test.createTemplateHelper(template)
	Expect(status).To(Equal(http.StatusCreated))
	Expect(templateID).NotTo(BeNil())
}

func (test TemplatesCRUD) GetTemplate(getTemplate params.Template) {
	templateID, _ := test.createTemplateHelper(getTemplate)
	statusCode, template := test.getTemplateHelper(templateID)

	Expect(statusCode).To(Equal(http.StatusOK))
	Expect(template).To(Equal(getTemplate))
}

func (test TemplatesCRUD) UpdateTemplate(updateTemplate params.Template) {
	templateID, _ := test.createTemplateHelper(updateTemplate)

	updateTemplate.Name = "New Name"
	updateTemplate.HTML = "<p>Brand new HTML</p>"
	updateTemplate.Subject = "lak;jsdfl;kajsdlf;"

	requestJSON, err := json.Marshal(updateTemplate)
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest("PUT", test.notificationsServer.TemplatePath(templateID), bytes.NewBuffer(requestJSON))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+test.clientToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	Expect(response.StatusCode).To(Equal(http.StatusNoContent))

	statusCode, template := test.getTemplateHelper(templateID)
	Expect(statusCode).To(Equal(http.StatusOK))
	Expect(template).To(Equal(updateTemplate))
}

func (test TemplatesCRUD) DeleteTemplate(deleteTemplate params.Template) {
	templateID, _ := test.createTemplateHelper(deleteTemplate)

	//delete existing template
	statusCode, body := test.deleteTemplateHelper(templateID)
	Expect(statusCode).To(Equal(http.StatusNoContent))
	Expect(body).To(BeEmpty())

	// get to verify 404
	statusCode, template := test.getTemplateHelper(templateID)
	Expect(statusCode).To(Equal(http.StatusNotFound))
	Expect(template).To(Equal(params.Template{}))

	// try to delete again (missing template) to verify 404
	statusCode, body = test.deleteTemplateHelper(templateID)
	Expect(statusCode).To(Equal(http.StatusNotFound))
	Expect(body).To(ContainSubstring("Not Found"))
}

func (test TemplatesCRUD) deleteTemplateHelper(templateID string) (int, []byte) {
	request, err := http.NewRequest("DELETE", test.notificationsServer.TemplatePath(templateID), bytes.NewBuffer([]byte{}))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+test.clientToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	return response.StatusCode, body
}

func (test TemplatesCRUD) getTemplateHelper(templateID string) (int, params.Template) {
	request, err := http.NewRequest("GET", test.notificationsServer.TemplatePath(templateID), bytes.NewBuffer([]byte{}))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+test.clientToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	if response.StatusCode != http.StatusOK {
		return response.StatusCode, params.Template{}
	}

	responseTemplate := params.Template{}
	err = json.Unmarshal(body, &responseTemplate)
	if err != nil {
		panic(err)
	}

	return response.StatusCode, responseTemplate
}

func (test TemplatesCRUD) createTemplateHelper(templateToCreate params.Template) (string, int) {
	jsonBody, err := json.Marshal(templateToCreate)
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest("POST", test.notificationsServer.TemplatesBasePath(), bytes.NewBuffer(jsonBody))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+test.clientToken.Access)
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
