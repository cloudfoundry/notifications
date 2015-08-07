package v2

import (
	"fmt"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/acceptance/v2/support"
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
	Metadata string
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
		token = GetClientTokenFor("my-client", "uaa")
	})

	It("can create a new template and retrieve it", func() {
		var createTemplate template

		By("creating a template", func() {
			status, response, err := client.Do("POST", "/templates", map[string]interface{}{
				"name":     "An interesting template",
				"text":     "template text",
				"html":     "template html",
				"subject":  "template subject",
				"metadata": `{"template": "metadata"}`,
			}, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusCreated))

			createTemplate.ID = response["id"].(string)
			createTemplate.Name = response["name"].(string)
			createTemplate.Text = response["text"].(string)
			createTemplate.Html = response["html"].(string)
			createTemplate.Subject = response["subject"].(string)
			createTemplate.Metadata = response["metadata"].(string)

			Expect(createTemplate.ID).NotTo(BeEmpty())
			Expect(createTemplate.Name).To(Equal("An interesting template"))
			Expect(createTemplate.Text).To(Equal("template text"))
			Expect(createTemplate.Html).To(Equal("template html"))
			Expect(createTemplate.Subject).To(Equal("template subject"))
			Expect(createTemplate.Metadata).To(MatchJSON(`{
				"template": "metadata"
			}`))
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
			getTemplate.Metadata = response["metadata"].(string)

			Expect(getTemplate.ID).To(Equal(createTemplate.ID))
			Expect(getTemplate.Name).To(Equal(createTemplate.Name))
			Expect(getTemplate.Html).To(Equal(createTemplate.Html))
			Expect(getTemplate.Text).To(Equal(createTemplate.Text))
			Expect(getTemplate.Subject).To(Equal(createTemplate.Subject))
			Expect(getTemplate.Metadata).To(Equal(createTemplate.Metadata))
		})
	})
})
