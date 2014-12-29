package acceptance

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/acceptance/support"
	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Template Metadata", func() {
	It("creates and updates a template with metadata", func() {
		var templateID string
		env := application.NewEnvironment()
		uaaClient := uaa.NewUAA("", env.UAAHost, "notifications-admin", "secret", "")
		clientToken, err := uaaClient.GetClientToken()
		if err != nil {
			panic(err)
		}
		client := support.NewClient(Servers.Notifications)

		By("creating a template with metadata", func() {
			var status int
			var err error

			status, templateID, err = client.Templates.Create(clientToken.Access, support.Template{
				Name:    "Star Wars",
				Subject: "Awesomeness",
				HTML:    "<p>Millenium Falcon</p>",
				Text:    "Millenium Falcon",
				Metadata: map[string]interface{}{
					"some_property": "some_value",
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusCreated))
			Expect(templateID).NotTo(BeNil())
		})

		By("verifying that the metadata was stored", func() {
			status, response, err := client.Templates.Get(clientToken.Access, templateID)

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(response.Metadata).To(Equal(map[string]interface{}{
				"some_property": "some_value",
			}))
		})

		By("updating the template metadata", func() {
			status, err := client.Templates.Update(clientToken.Access, templateID, support.Template{
				Name:    "Star Wars",
				Subject: "Awesomeness",
				HTML:    "<p>Millenium Falcon</p>",
				Text:    "Millenium Falcon",
				Metadata: map[string]interface{}{
					"hello": true,
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNoContent))
		})

		By("verifying that the metadata was updated", func() {
			status, response, err := client.Templates.Get(clientToken.Access, templateID)

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(response.Metadata).To(Equal(map[string]interface{}{
				"hello": true,
			}))
		})
	})
})
