package v1

import (
	"net/http"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/cloudfoundry-incubator/notifications/v1/acceptance/support"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Send a notification to an email", func() {
	It("sends a single notification to an email", func() {
		var templateID string
		var response support.NotifyResponse
		clientID := "notifications-sender"
		clientToken := GetClientTokenFor(clientID)
		client := support.NewClient(Servers.Notifications.URL())

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
				HTML:    "<p>Enterprise</p>{{.HTML}}<h1>{{.Endorsement}}</h1>{{.Domain}}",
				Text:    "Enterprise\n{{.Text}}\n{{.Endorsement}}",
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
			status, responses, err := client.Notify.Email(clientToken.Access, "John User <user@example.com>", support.Notify{
				HTML:    "<header>this is an acceptance test</header>",
				Text:    "some text for the email",
				Subject: "my-special-subject",
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))

			Expect(responses).To(HaveLen(1))
			response = responses[0]
			Expect(response.Status).To(Equal("queued"))
			Expect(response.Recipient).To(Equal("user@example.com"))
			Expect(response.VCAPRequestID).To(Equal("some-totally-fake-vcap-request-id"))
			Expect(GUIDRegex.MatchString(response.NotificationID)).To(BeTrue())
		})

		By("verifying the message was sent", func() {
			Eventually(func() int {
				return len(Servers.SMTP.Deliveries)
			}, 10*time.Second).Should(Equal(1))
			delivery := Servers.SMTP.Deliveries[0]

			env, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())

			Expect(delivery.Sender).To(Equal(env.Sender))
			Expect(delivery.Recipients).To(Equal([]string{"user@example.com"}))

			data := strings.Split(string(delivery.Data), "\n")
			Expect(data).To(ContainElement("X-CF-Client-ID: notifications-sender"))
			Expect(data).To(ContainElement("X-CF-Notification-ID: " + response.NotificationID))
			Expect(data).To(ContainElement(ContainSubstring("X-CF-Notification-Timestamp: ")))

			var timestampString string
			prefix := "X-CF-Notification-Timestamp: "
			for _, line := range data {
				if !strings.Contains(line, prefix) {
					continue
				}
				timestampString = strings.TrimPrefix(line, prefix)
			}
			Expect(timestampString).NotTo(BeEmpty())

			timestampDate, err := time.Parse(time.RFC3339, timestampString)
			Expect(err).NotTo(HaveOccurred())
			Expect(timestampDate).To(BeTemporally("~", time.Now(), 15*time.Second))

			Expect(data).To(ContainElement("Subject: Boldness my-special-subject"))
			Expect(data).To(ContainElement("\t\t<p>Enterprise</p><header>this is an acceptance test</header><h1>This messa="))
			Expect(data).To(ContainElement("ge was sent directly to your email address.</h1>localhost"))
			Expect(data).To(ContainElement("Enterprise"))
			Expect(data).To(ContainElement("some text for the email"))
			Expect(data).To(ContainElement("This message was sent directly to your email address."))
		})

		By("confirming that the client notifications list remains unaffected", func() {
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
