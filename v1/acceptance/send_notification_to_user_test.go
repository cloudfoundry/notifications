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

var _ = Describe("Send a notification to a user", func() {
	It("sends a single notification email to a user", func() {
		var (
			templateID  string
			response    support.NotifyResponse
			clientID    = "notifications-sender"
			clientToken = GetClientTokenFor(clientID)
			client      = support.NewClient(Servers.Notifications.URL())
			userID      = "user-123"
		)

		env, err := application.NewEnvironment()
		Expect(err).NotTo(HaveOccurred())

		By("registering a notification", func() {
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
				Name:    "Star Wars",
				Subject: "Awesomeness {{.Subject}}",
				HTML:    "<p>Millenium Falcon</p>{{.HTML}}<b>{{.Endorsement}}</b>",
				Text:    "Millenium Falcon\n{{.Text}}\n{{.Endorsement}}",
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

		By("sending a notifications to a user", func() {
			status, responses, err := client.Notify.User(clientToken.Access, userID, support.Notify{
				KindID:  "acceptance-test",
				HTML:    "<p>this is an acceptance%40test</p>",
				Text:    "hello from the acceptance test",
				Subject: "my-special-subject",
				ReplyTo: "males@example.com",
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))

			Expect(responses).To(HaveLen(1))
			response = responses[0]
			Expect(response.Status).To(Equal("queued"))
			Expect(response.Recipient).To(Equal(userID))
			Expect(GUIDRegex.MatchString(response.NotificationID)).To(BeTrue())
			Expect(response.VCAPRequestID).To(Equal("some-totally-fake-vcap-request-id"))
		})

		By("verifying that the message was sent", func() {
			Eventually(func() int {
				return len(Servers.SMTP.Deliveries)
			}, 10*time.Second).Should(Equal(1))
			delivery := Servers.SMTP.Deliveries[0]

			Expect(delivery.Sender).To(Equal(env.Sender))
			Expect(delivery.Recipients).To(Equal([]string{"user-123@example.com"}))

			data := strings.Split(string(delivery.Data), "\n")
			Expect(data).To(ContainElement("X-CF-Client-ID: notifications-sender"))
			Expect(data).To(ContainElement("X-CF-Notification-ID: " + response.NotificationID))
			Expect(data).To(ContainElement(HavePrefix("X-CF-Notification-Request-Received: ")))
			Expect(data).To(ContainElement("Reply-To: males@example.com"))
			Expect(data).To(ContainElement("Subject: Awesomeness my-special-subject"))
			Expect(data).To(ContainElement("\t\t<p>Millenium Falcon</p><p>this is an acceptance%40test</p><b>This message ="))
			Expect(data).To(ContainElement("was sent directly to you.</b>"))
			Expect(data).To(ContainElement("hello from the acceptance test"))
			Expect(data).To(ContainElement("This message was sent directly to you."))

			for _, line := range data {
				if strings.HasPrefix(line, "X-CF-Notification-Request-Received: ") {
					reqRecTime, _ := time.Parse(time.RFC3339Nano, strings.Split(line, ": ")[1])
					Expect(reqRecTime).To(BeTemporally("~", time.Now(), 5*time.Minute))
				}
			}
		})
	})

	Context("when the workers fail", func() {
		AfterEach(func() {
			Servers.Notifications.ResetDatabase()
		})

		It("retries deliveries when they fail to be sent", func() {
			var (
				templateID  string
				response    support.NotifyResponse
				clientID    = "notifications-sender"
				clientToken = GetClientTokenFor(clientID)
				client      = support.NewClient(Servers.Notifications.URL())
				userID      = "user-malformed-email"
			)

			By("registering a notification", func() {
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
					Name:    "Star Wars",
					Subject: "Awesomeness {{.Subject}}",
					HTML:    "<p>Millenium Falcon</p>{{.HTML}}<b>{{.Endorsement}}</b>",
					Text:    "Millenium Falcon\n{{.Text}}\n{{.Endorsement}}",
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

			By("sending a notifications to a user", func() {
				status, responses, err := client.Notify.User(clientToken.Access, userID, support.Notify{
					KindID:  "acceptance-test",
					HTML:    "<p>this is an acceptance%40test</p>",
					Text:    "hello from the acceptance test",
					Subject: "my-special-subject",
					ReplyTo: "males@example.com",
				})

				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusOK))

				Expect(responses).To(HaveLen(1))
				response = responses[0]

				Expect(response.Status).To(Equal("queued"))
				Expect(response.Recipient).To(Equal(userID))
				Expect(GUIDRegex.MatchString(response.NotificationID)).To(BeTrue())
				Expect(response.VCAPRequestID).To(Equal("some-totally-fake-vcap-request-id"))
			})

			By("verifying that the message never gets sent", func() {
				Consistently(func() int {
					return len(Servers.SMTP.Deliveries)
				}).Should(Equal(0))
			})
		})
	})
})
