package v1

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/v1/acceptance/support"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Templates CRUD", func() {
	var (
		templates   []support.Template
		client      *support.Client
		clientToken uaa.Token
	)

	BeforeEach(func() {
		clientToken = GetClientTokenFor("notifications-admin")
		client = support.NewClient(Servers.Notifications.URL())

		templates = []support.Template{
			{
				Name:     "Star Wars",
				Subject:  "Awesomeness",
				HTML:     "<p>Millenium Falcon</p>",
				Text:     "Millenium Falcon",
				Metadata: make(map[string]interface{}),
			},
			{
				Name:    "Big Hero 6",
				Subject: "Heroes",
				HTML:    "<h1>Robots!</h1>",
				Text:    "Robots!",
			},
			{
				Name:     "Blah",
				Subject:  "More Blah",
				HTML:     "<h1>This is blahblah</h1>",
				Text:     "Blah even more",
				Metadata: make(map[string]interface{}),
			},
			{
				Name:     "Hungry Play",
				Subject:  "Dystopian",
				HTML:     "<h1>Sad</h1>",
				Text:     "Run!!",
				Metadata: make(map[string]interface{}),
			},
		}
	})

	It("allows a user to create a new template", func() {
		status, templateID, err := client.Templates.Create(clientToken.Access, templates[0])
		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(http.StatusCreated))
		Expect(templateID).NotTo(BeNil())
	})

	It("errors when a bad template variable is used", func() {
		status, _, err := client.Templates.Create(clientToken.Access, support.Template{
			Name:    "Bad Template",
			Subject: "This is a bad template",
			HTML:    "<p>{{some-bad-arg}}</p>",
			Text:    "hello {{some-bad-arg}}",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(422))
	})

	It("allows a user to retrieve a template", func() {
		var templateID string

		By("creating a template without metadata field set", func() {
			var err error
			_, templateID, err = client.Templates.Create(clientToken.Access, templates[1])
			Expect(err).NotTo(HaveOccurred())
		})

		By("verifying that the template can be retrieved, and has metadata set to {}", func() {
			status, template, err := client.Templates.Get(clientToken.Access, templateID)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(template).To(Equal(support.Template{
				Name:     "Big Hero 6",
				Subject:  "Heroes",
				HTML:     "<h1>Robots!</h1>",
				Text:     "Robots!",
				Metadata: make(map[string]interface{}),
			}))
		})
	})

	It("allows a user to update an existing template", func() {
		var templateID string

		By("creating a template", func() {
			var err error
			_, templateID, err = client.Templates.Create(clientToken.Access, templates[2])
			Expect(err).NotTo(HaveOccurred())
		})

		By("updating the template data", func() {
			templates[2].Name = "New Name"
			templates[2].HTML = "<p>Brand new HTML</p>"
			templates[2].Subject = "lak;jsdfl;kajsdlf;"

			status, err := client.Templates.Update(clientToken.Access, templateID, templates[2])
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNoContent))
		})

		By("verifying that the template was updated", func() {
			status, actualTemplate, err := client.Templates.Get(clientToken.Access, templateID)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(actualTemplate).To(Equal(templates[2]))
		})
	})

	It("allows a user to delete a template", func() {
		var templateID string

		By("creating a template", func() {
			var err error
			_, templateID, err = client.Templates.Create(clientToken.Access, templates[3])
			Expect(err).NotTo(HaveOccurred())
		})

		By("deleting the template", func() {
			status, err := client.Templates.Delete(clientToken.Access, templateID)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNoContent))
		})

		By("verifying that the template no longer exists", func() {
			status, _, err := client.Templates.Get(clientToken.Access, templateID)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNotFound))
		})

		By("verifying that the template cannot be deleted again", func() {
			status, err := client.Templates.Delete(clientToken.Access, templateID)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNotFound))
		})
	})

	It("allows a user to list all templates", func() {
		templatesList := []support.TemplateListItem{}

		By("creating several templates", func() {
			for _, template := range templates {
				status, templateID, err := client.Templates.Create(clientToken.Access, template)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusCreated))

				templatesList = append(templatesList, support.TemplateListItem{
					ID:   templateID,
					Name: template.Name,
				})
			}
		})

		By("verifying that the created templates are listed", func() {
			status, actualTemplates, err := client.Templates.List(clientToken.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(200))
			Expect(actualTemplates).To(HaveLen(4))
			for _, template := range templatesList {
				Expect(actualTemplates).To(ContainElement(template))
			}
		})
	})
})
