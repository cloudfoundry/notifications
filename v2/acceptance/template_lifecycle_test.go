package acceptance

import (
	"fmt"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/v2/acceptance/support"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

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

	It("can create a new template, retrieve it and delete it again", func() {
		var templateID string

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

			templateID = response["id"].(string)

			Expect(templateID).NotTo(BeEmpty())
			Expect(response["name"]).To(Equal("An interesting template"))
			Expect(response["text"]).To(Equal("template text"))
			Expect(response["html"]).To(Equal("template html"))
			Expect(response["subject"]).To(Equal("template subject"))
			Expect(response["metadata"]).To(Equal(map[string]interface{}{
				"template": "metadata",
			}))
		})

		By("getting the template", func() {
			status, response, err := client.Do("GET", fmt.Sprintf("/templates/%s", templateID), nil, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))

			Expect(response["id"]).To(Equal(templateID))
			Expect(response["name"]).To(Equal("An interesting template"))
			Expect(response["text"]).To(Equal("template text"))
			Expect(response["html"]).To(Equal("template html"))
			Expect(response["subject"]).To(Equal("template subject"))
			Expect(response["metadata"]).To(Equal(map[string]interface{}{
				"template": "metadata",
			}))
		})

		By("deleting the template", func() {
			status, _, err := client.Do("DELETE", fmt.Sprintf("/templates/%s", templateID), nil, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNoContent))
		})

		By("failing to get the deleted template", func() {
			status, _, err := client.Do("GET", fmt.Sprintf("/templates/%s", templateID), nil, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNotFound))
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

		It("returns a 404 when the template to delete does not exist", func() {
			status, response, err := client.Do("DELETE", "/templates/missing-template-id", nil, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNotFound))
			Expect(response["errors"]).To(ContainElement("Template with id \"missing-template-id\" could not be found"))
		})
	})
})
