package v1

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/v1/acceptance/support"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Assign Templates", func() {
	var (
		templateID     string
		notificationID string
		clientID       string
		clientToken    uaa.Token
		client         *support.Client
	)

	BeforeEach(func() {
		notificationID = "acceptance-test"
		clientID = "notifications-admin"
		clientToken = GetClientTokenFor(clientID)
		client = support.NewClient(Servers.Notifications.URL())
	})

	It("Creates a template and then assigns it", func() {
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

		By("confirming that the client and notification are listed in the assignments", func() {
			status, associations, err := client.Templates.Associations(clientToken.Access, templateID)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(associations).To(HaveLen(2))
			Expect(associations).To(ContainElement(support.TemplateAssociation{
				ClientID: clientID,
			}))
			Expect(associations).To(ContainElement(support.TemplateAssociation{
				ClientID:       clientID,
				NotificationID: notificationID,
			}))
		})

		By("resetting the template assignment", func() {
			status, err := client.Templates.AssignToNotification(clientToken.Access, clientID, notificationID, "")
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNoContent))
		})

		By("confirming that the notifications has the default template assigned", func() {
			status, notifications, err := client.Notifications.List(clientToken.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(notifications).To(HaveLen(1))

			clientNotifications := notifications[clientID].Notifications
			Expect(clientNotifications).To(HaveLen(1))
			Expect(clientNotifications[notificationID].Template).To(Equal("default"))
		})
	})

	It("does not know about non-existent templates", func() {
		status, _, err := client.Templates.Associations(clientToken.Access, "random-template-id")
		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(http.StatusNotFound))
	})
})
