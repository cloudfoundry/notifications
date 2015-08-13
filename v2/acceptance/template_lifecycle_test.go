package v2

import (
	"fmt"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/v2/acceptance/support"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type template struct {
	ID       string
	Name     string
	Text     string
	Html     string
	Subject  string
	Metadata map[string]interface{}
}

var _ = Describe("Template lifecycle", func() {
	var (
		client *support.Client
		token  uaa.Token
	)

	BeforeEach(func() {
		client = support.NewClient(support.Config{
			Host:  Servers.Notifications.URL(),
			Trace: Trace,
		})
		token = GetClientTokenFor("my-client")
	})

	It("can create a new template and retrieve it", func() {
		var createTemplate template

		By("creating a template", func() {
			status, response, err := client.Do("POST", "/templates", map[string]interface{}{
				"name":    "An interesting template",
				"text":    "template text",
				"html":    "template html",
				"subject": "template subject",
				"metadata": map[string]interface{}{
					"template": "metadata",
				},
			}, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusCreated))

			createTemplate.ID = response["id"].(string)
			createTemplate.Name = response["name"].(string)
			createTemplate.Text = response["text"].(string)
			createTemplate.Html = response["html"].(string)
			createTemplate.Subject = response["subject"].(string)
			createTemplate.Metadata = response["metadata"].(map[string]interface{})

			Expect(createTemplate.ID).NotTo(BeEmpty())
			Expect(createTemplate.Name).To(Equal("An interesting template"))
			Expect(createTemplate.Text).To(Equal("template text"))
			Expect(createTemplate.Html).To(Equal("template html"))
			Expect(createTemplate.Subject).To(Equal("template subject"))
			Expect(createTemplate.Metadata).To(Equal(map[string]interface{}{
				"template": "metadata",
			}))
		})

		By("getting a template", func() {
			status, response, err := client.Do("GET", fmt.Sprintf("/templates/%s", createTemplate.ID), nil, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))

			var getTemplate template
			getTemplate.ID = response["id"].(string)
			getTemplate.Name = response["name"].(string)
			getTemplate.Text = response["text"].(string)
			getTemplate.Html = response["html"].(string)
			getTemplate.Subject = response["subject"].(string)
			getTemplate.Metadata = response["metadata"].(map[string]interface{})

			Expect(getTemplate.ID).To(Equal(createTemplate.ID))
			Expect(getTemplate.Name).To(Equal(createTemplate.Name))
			Expect(getTemplate.Html).To(Equal(createTemplate.Html))
			Expect(getTemplate.Text).To(Equal(createTemplate.Text))
			Expect(getTemplate.Subject).To(Equal(createTemplate.Subject))
			Expect(getTemplate.Metadata).To(Equal(createTemplate.Metadata))
		})
	})

	Context("failure states", func() {
		It("returns a 409 with the correct error message when a template already exists", func() {
			status, response, err := client.Do("POST", "/templates", map[string]interface{}{
				"name":    "An interesting template",
				"text":    "template text",
				"html":    "template html",
				"subject": "template subject",
				"metadata": map[string]interface{}{
					"template": "metadata",
				},
			}, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusCreated))

			status, response, err = client.Do("POST", "/templates", map[string]interface{}{
				"name":    "An interesting template",
				"text":    "template text",
				"html":    "template html",
				"subject": "template subject",
				"metadata": map[string]interface{}{
					"template": "metadata",
				},
			}, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusConflict))
			Expect(response["errors"]).To(ContainElement("Template with name \"An interesting template\" already exists"))
		})

		It("returns a 404 when the template cannot be retrieved", func() {
			status, response, err := client.Do("GET", "/templates/missing-template-id", nil, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNotFound))
			Expect(response["errors"]).To(ContainElement("Template with id \"missing-template-id\" could not be found"))
		})

		It("returns a 404 when the template belongs to a different client", func() {
			var templateID string

			By("creating a template for one client", func() {
				status, response, err := client.Do("POST", "/templates", map[string]interface{}{
					"name":    "An interesting template",
					"text":    "template text",
					"html":    "template html",
					"subject": "template subject",
					"metadata": map[string]interface{}{
						"template": "metadata",
					},
				}, token.Access)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusCreated))

				templateID = response["id"].(string)
			})

			By("attempting to access the created template as another client", func() {
				token := GetClientTokenFor("other-client")
				status, response, err := client.Do("GET", fmt.Sprintf("/templates/%s", templateID), nil, token.Access)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusNotFound))
				Expect(response["errors"]).To(ContainElement(fmt.Sprintf("Template with id %q could not be found", templateID)))
			})
		})
	})
})
