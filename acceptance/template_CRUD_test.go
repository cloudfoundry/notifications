package acceptance

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/acceptance/servers"
	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/cloudfoundry-incubator/notifications/web/params"
	"github.com/cloudfoundry-incubator/notifications/web/services"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Templates CRUD", func() {
	var t TemplatesCRUD
	var testTemplates []params.Template

	BeforeEach(func() {
		env := application.NewEnvironment()
		uaaClient := uaa.NewUAA("", env.UAAHost, "notifications-admin", "secret", "")
		clientToken, err := uaaClient.GetClientToken()
		if err != nil {
			panic(err)
		}

		t = TemplatesCRUD{
			notificationsServer: Servers.Notifications,
			clientToken:         clientToken,
		}

		testTemplates = []params.Template{
			params.Template{
				Name:     "Star Wars",
				Subject:  "Awesomeness",
				HTML:     "<p>Millenium Falcon</p>",
				Text:     "Millenium Falcon",
				Metadata: make(map[string]interface{}),
			},
			params.Template{
				Name:     "Big Hero 6",
				Subject:  "Heroes",
				HTML:     "<h1>Robots!</h1>",
				Text:     "Robots!",
				Metadata: make(map[string]interface{}),
			},
			params.Template{
				Name:     "Blah",
				Subject:  "More Blah",
				HTML:     "<h1>This is blahblah</h1>",
				Text:     "Blah even more",
				Metadata: make(map[string]interface{}),
			},
			params.Template{
				Name:     "Hungry Play",
				Subject:  "Dystopian",
				HTML:     "<h1>Sad</h1>",
				Text:     "Run!!",
				Metadata: make(map[string]interface{}),
			},
		}
	})

	It("allows a user to create a new template", func() {
		status, templateID := t.createTemplateHelper(testTemplates[0])
		Expect(status).To(Equal(http.StatusCreated))
		Expect(templateID).NotTo(BeNil())
	})

	It("allows a user to retrieve a template", func() {
		_, templateID := t.createTemplateHelper(testTemplates[1])
		statusCode, template := t.getTemplateHelper(templateID)

		Expect(statusCode).To(Equal(http.StatusOK))
		Expect(template).To(Equal(testTemplates[1]))
	})

	It("allows a user to update an existing template", func() {
		template := testTemplates[2]
		_, templateID := t.createTemplateHelper(template)

		template.Name = "New Name"
		template.HTML = "<p>Brand new HTML</p>"
		template.Subject = "lak;jsdfl;kajsdlf;"

		requestJSON, err := json.Marshal(template)
		if err != nil {
			panic(err)
		}

		request, err := http.NewRequest("PUT", t.notificationsServer.TemplatePath(templateID), bytes.NewBuffer(requestJSON))
		if err != nil {
			panic(err)
		}

		request.Header.Set("Authorization", "Bearer "+t.clientToken.Access)

		response, err := http.DefaultClient.Do(request)
		if err != nil {
			panic(err)
		}

		Expect(response.StatusCode).To(Equal(http.StatusNoContent))

		statusCode, template := t.getTemplateHelper(templateID)
		Expect(statusCode).To(Equal(http.StatusOK))
		Expect(template).To(Equal(template))
	})

	It("allows a user to delete a template", func() {
		deleteTemplate := testTemplates[3]
		_, templateID := t.createTemplateHelper(deleteTemplate)

		//delete existing template
		statusCode, body := t.deleteTemplateHelper(templateID)
		Expect(statusCode).To(Equal(http.StatusNoContent))
		Expect(bufio.NewReader(body).Buffered()).To(Equal(0))

		// get to verify 404
		statusCode, template := t.getTemplateHelper(templateID)
		Expect(statusCode).To(Equal(http.StatusNotFound))
		Expect(template).To(Equal(params.Template{}))

		// try to delete again (missing template) to verify 404
		statusCode, body = t.deleteTemplateHelper(templateID)
		Expect(statusCode).To(Equal(http.StatusNotFound))
		buffer := bytes.NewBuffer([]byte{})
		_, err := buffer.ReadFrom(body)
		if err != nil {
			panic(err)
		}
		Expect(buffer).To(ContainSubstring("Not Found"))
	})

	It("allows a user to list all templates", func() {
		//create a bunch of templates
		templateMetadata := map[string]services.TemplateMetadata{}
		for _, fullTemplate := range testTemplates {
			statusCode, templateID := t.createTemplateHelper(fullTemplate)
			if statusCode != http.StatusCreated {
				panic("ListTemplates failed to create test Templates")
			}
			templateMetadata[templateID] = services.TemplateMetadata{Name: fullTemplate.Name}
		}

		//call Get /templates
		request, err := http.NewRequest("GET", t.notificationsServer.TemplatesBasePath(), bytes.NewBuffer([]byte{}))
		request.Header.Set("Authorization", "Bearer "+t.clientToken.Access)
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			panic(err)
		}

		Expect(response.StatusCode).To(Equal(200))

		var templatesListResponse map[string]services.TemplateMetadata
		err = json.NewDecoder(response.Body).Decode(&templatesListResponse)
		if err != nil {
			panic(err)
		}

		Expect(templatesListResponse).To(Equal(templateMetadata))
	})
})

type TemplatesCRUD struct {
	notificationsServer servers.Notifications
	clientToken         uaa.Token
}

func (test TemplatesCRUD) deleteTemplateHelper(templateID string) (int, io.Reader) {
	request, err := http.NewRequest("DELETE", test.notificationsServer.TemplatePath(templateID), bytes.NewBuffer([]byte{}))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", "Bearer "+test.clientToken.Access)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	return response.StatusCode, response.Body
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

	if response.StatusCode != http.StatusOK {
		return response.StatusCode, params.Template{}
	}

	responseTemplate := params.Template{}
	err = json.NewDecoder(response.Body).Decode(&responseTemplate)
	if err != nil {
		panic(err)
	}

	return response.StatusCode, responseTemplate
}

func (test TemplatesCRUD) createTemplateHelper(templateToCreate params.Template) (int, string) {
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

	var JSON struct {
		TemplateID string `json:"template_id"`
	}

	err = json.NewDecoder(response.Body).Decode(&JSON)
	if err != nil {
		panic(err)
	}

	return response.StatusCode, JSON.TemplateID
}
