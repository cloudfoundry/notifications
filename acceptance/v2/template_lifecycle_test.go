package v2

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/acceptance/v2/support"
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
		token = GetClientTokenFor("my-client", "uaa")
	})

	It("can create a new template", func() {
		var template struct {
			ID       string
			Name     string
			Text     string
			Html     string
			Subject  string
			Metadata string
		}

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

			template.ID = response["id"].(string)
			template.Name = response["name"].(string)
			template.Text = response["text"].(string)
			template.Html = response["html"].(string)
			template.Subject = response["subject"].(string)
			template.Metadata = response["metadata"].(string)

			Expect(template.ID).NotTo(BeEmpty())
			Expect(template.Name).To(Equal("An interesting template"))
			Expect(template.Text).To(Equal("template text"))
			Expect(template.Html).To(Equal("template html"))
			Expect(template.Subject).To(Equal("template subject"))
			Expect(template.Metadata).To(MatchJSON(`{
				"template": "metadata"
			}`))
		})
	})
})
