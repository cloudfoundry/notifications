package v1

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/v1/acceptance/support"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Updating A Notification", func() {
	var (
		client         *support.Client
		clientToken    uaa.Token
		notificationID string
		clientID       string
	)

	BeforeEach(func() {
		notificationID = "acceptance-test"
		clientID = "notifications-admin"
		clientToken = GetClientTokenFor(clientID)
		client = support.NewClient(Servers.Notifications.URL())
	})

	It("can update notifications", func() {
		By("registering a notification", func() {
			code, err := client.Notifications.Register(clientToken.Access, support.RegisterClient{
				SourceName: "Notifications Sender",
				Notifications: map[string]support.RegisterNotification{
					notificationID: {
						Description: "Acceptance Test",
						Critical:    false,
					},
				},
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(http.StatusNoContent))
		})

		By("updating a notification", func() {
			updatedNotification := support.Notification{
				Description: "Acceptance Test With Modified Description",
				Template:    "New Template",
				Critical:    true,
			}

			status, err := client.Notifications.Update(clientToken.Access, clientID, notificationID, updatedNotification)

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNoContent))
		})

		By("verifying the notification was updated", func() {
			status, notifications, err := client.Notifications.List(clientToken.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(notifications).To(HaveLen(1))

			clientNotifications := notifications[clientID].Notifications
			Expect(clientNotifications).To(HaveLen(1))
			Expect(clientNotifications[notificationID].Description).To(Equal("Acceptance Test With Modified Description"))
			Expect(clientNotifications[notificationID].Template).To(Equal("New Template"))
			Expect(clientNotifications[notificationID].Critical).To(Equal(true))
		})
	})
})
