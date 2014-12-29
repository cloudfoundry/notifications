package acceptance

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/acceptance/support"
	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Default Template", func() {
	var client *support.Client
	var env application.Environment
	var clientToken uaa.Token

	BeforeEach(func() {
		var err error
		clientID := "notifications-admin"
		env = application.NewEnvironment()
		uaaClient := uaa.NewUAA("", env.UAAHost, clientID, "secret", "")
		clientToken, err = uaaClient.GetClientToken()
		if err != nil {
			panic(err)
		}
		client = support.NewClient(Servers.Notifications)
	})

	It("can retrieve the default template", func() {
		status, template, err := client.Templates.Default.Get(clientToken.Access)
		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(http.StatusOK))
		Expect(template).To(Equal(support.Template{
			Name:     "Default Template",
			Subject:  "CF Notification: {{.Subject}}",
			HTML:     "{{.HTML}}",
			Text:     "{{.Text}}",
			Metadata: map[string]interface{}{},
		}))
	})

	It("can edit the default template", func() {
		By("editing the default template", func() {
			status, err := client.Templates.Default.Update(clientToken.Access, support.Template{
				Name:    "A Whole New Template",
				Subject: "Updated: {{.Subject}}",
				HTML:    "<h1>Updated!!!</h1>",
				Text:    "Updated!!!",
				Metadata: map[string]interface{}{
					"smurf": "favorite",
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNoContent))
		})

		By("verifying that the default template was updated", func() {
			status, template, err := client.Templates.Default.Get(clientToken.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(template).To(Equal(support.Template{
				Name:    "A Whole New Template",
				Subject: "Updated: {{.Subject}}",
				HTML:    "<h1>Updated!!!</h1>",
				Text:    "Updated!!!",
				Metadata: map[string]interface{}{
					"smurf": "favorite",
				},
			}))
		})

		By("restarting the notifications service", func() {
			Servers.Notifications.Restart()
		})

		By("verifying that the default template still displays the overridden values", func() {
			status, template, err := client.Templates.Default.Get(clientToken.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(template).To(Equal(support.Template{
				Name:    "A Whole New Template",
				Subject: "Updated: {{.Subject}}",
				HTML:    "<h1>Updated!!!</h1>",
				Text:    "Updated!!!",
				Metadata: map[string]interface{}{
					"smurf": "favorite",
				},
			}))
		})
	})
})
