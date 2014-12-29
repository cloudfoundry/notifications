package acceptance

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/acceptance/support"
	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Assign Templates", func() {
	It("Creates a template and then assigns it", func() {
		var templateID string
		notificationID := "acceptance-test"
		clientID := "notifications-admin"
		env := application.NewEnvironment()
		uaaClient := uaa.NewUAA("", env.UAAHost, clientID, "secret", "")
		clientToken, err := uaaClient.GetClientToken()
		if err != nil {
			panic(err)
		}

		client := support.NewClient(Servers.Notifications)

		By("registering a notification", func() {
			code, err := client.Notifications.Register(clientToken.Access, support.RegisterClient{
				SourceName: "Notifications Sender",
				Notifications: map[string]support.RegisterNotification{
					notificationID: {
						Description: "Acceptance Test",
						Critical:    true,
					},
				},
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(http.StatusNoContent))
		})

		By("creating a template", func() {
			var status int
			var err error

			status, templateID, err = client.Templates.Create(clientToken.Access, support.Template{
				Name:    "Star Wars",
				Subject: "Awesomeness",
				HTML:    "<p>Millenium Falcon</p>",
				Text:    "Millenium Falcon",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusCreated))
			Expect(templateID).NotTo(BeNil())
		})

		By("assigning the template to a client", func() {
			status, err := client.Templates.AssignToClient(clientToken.Access, clientID, templateID)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNoContent))
		})

		By("confirming that the client has the assigned template", func() {
			status, notifications, err := client.Notifications.List(clientToken.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(notifications).To(HaveLen(1))
			Expect(notifications[clientID].Template).To(Equal(templateID))
		})

		By("assigning the template to a notification", func() {
			status, err := client.Templates.AssignToNotification(clientToken.Access, clientID, notificationID, templateID)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNoContent))
		})

		By("confirming that the notifications has the assigned template", func() {
			status, notifications, err := client.Notifications.List(clientToken.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(notifications).To(HaveLen(1))

			clientNotifications := notifications[clientID].Notifications
			Expect(clientNotifications).To(HaveLen(1))
			Expect(clientNotifications[notificationID].Template).To(Equal(templateID))
		})
	})
})
