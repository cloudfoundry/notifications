package v1

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/v1/acceptance/support"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Template Metadata", func() {
	It("creates and updates a template with metadata", func() {
		var templateID string
		clientToken := GetClientTokenFor("notifications-admin")
		client := support.NewClient(Servers.Notifications.URL())

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

		By("setting a template without metadata field set", func() {
			status, err := client.Templates.Update(clientToken.Access, templateID, support.Template{
				Name:    "Flashy Wars",
				Subject: "ness",
				HTML:    "<p>Alcon</p>",
				Text:    "Alcon",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNoContent))
		})

		By("verifying the metadata is set to {}", func() {
			status, response, err := client.Templates.Get(clientToken.Access, templateID)

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(response.Metadata).To(Equal(map[string]interface{}{}))
		})
	})
})
