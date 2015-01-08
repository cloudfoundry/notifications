package acceptance

import (
	"net/http"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/notifications/acceptance/support"
	"github.com/cloudfoundry-incubator/notifications/application"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Send a notification to an email", func() {
	It("sends a single notification to an email", func() {
		var templateID string
		var response support.NotifyResponse
		clientID := "notifications-sender"
		clientToken := GetClientTokenFor(clientID)
		client := support.NewClient(Servers.Notifications)

		By("registering a notifications", func() {
			code, err := client.Notifications.Register(clientToken.Access, support.RegisterClient{
				SourceName: "Notifications Sender",
				Notifications: map[string]support.RegisterNotification{
					"acceptance-test": {
						Description: "Acceptance Test",
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(http.StatusNoContent))
		})

		By("creating a new template", func() {
			var status int
			var err error
			status, templateID, err = client.Templates.Create(clientToken.Access, support.Template{
				Name:    "Star Trek",
				Subject: "Boldness {{.Subject}}",
				HTML:    "<p>Enterprise</p>{{.HTML}}",
				Text:    "Enterprise\n{{.Text}}",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusCreated))
			Expect(templateID).NotTo(Equal(""))
		})

		By("assigning the template to a client", func() {
			status, err := client.Templates.AssignToClient(clientToken.Access, clientID, templateID)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNoContent))
		})

		By("sending a notification to an email address", func() {
			status, responses, err := client.Notify.Email(clientToken.Access, support.Notify{
				HTML:    "<header>this is an acceptance test</header>",
				Subject: "my-special-subject",
				To:      "John User <user@example.com>",
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))

			Expect(responses).To(HaveLen(1))
			response = responses[0]
			Expect(response.Status).To(Equal("queued"))
			Expect(response.Recipient).To(Equal("user@example.com"))
			Expect(GUIDRegex.MatchString(response.NotificationID)).To(BeTrue())
		})

		By("verifying the message was sent", func() {
			Eventually(func() int {
				return len(Servers.SMTP.Deliveries)
			}, 1*time.Second).Should(Equal(1))
			delivery := Servers.SMTP.Deliveries[0]

			env := application.NewEnvironment()
			Expect(delivery.Sender).To(Equal(env.Sender))
			Expect(delivery.Recipients).To(Equal([]string{"user@example.com"}))

			data := strings.Split(string(delivery.Data), "\n")
			Expect(data).To(ContainElement("X-CF-Client-ID: notifications-sender"))
			Expect(data).To(ContainElement("X-CF-Notification-ID: " + response.NotificationID))
			Expect(data).To(ContainElement("Subject: Boldness my-special-subject"))
			Expect(data).To(ContainElement("        <p>Enterprise</p><header>this is an acceptance test</header>"))
		})

		By("confirming that the client notificatins list remains unaffected", func() {
			status, list, err := client.Notifications.List(GetClientTokenFor("notifications-sender").Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(list).To(HaveLen(1))

			notificationsSender := list[clientID]
			Expect(notificationsSender.Name).To(Equal("Notifications Sender"))
			Expect(notificationsSender.Template).To(Equal(templateID))
			Expect(notificationsSender.Notifications).To(HaveLen(1))

			acceptanceNotification := notificationsSender.Notifications["acceptance-test"]
			Expect(acceptanceNotification.Description).To(Equal("Acceptance Test"))
			Expect(acceptanceNotification.Template).To(Equal("default"))
			Expect(acceptanceNotification.Critical).To(BeFalse())
		})
	})
})
