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

var _ = Describe("Sending notifications to all users in an organization", func() {
	It("sends a notification to each user in an organization", func() {
		var templateID string
		indexedResponses := map[string]support.NotifyResponse{}
		clientID := "notifications-sender"
		clientToken := GetClientTokenFor(clientID)
		client := support.NewClient(Servers.Notifications.URL())

		By("registering a notification", func() {
			status, err := client.Notifications.Register(clientToken.Access, support.RegisterClient{
				SourceName: "Notifications Sender",
				Notifications: map[string]support.RegisterNotification{
					"organization-test": {
						Description: "Organization Test",
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNoContent))
		})

		By("creating a template", func() {
			var status int
			var err error
			status, templateID, err = client.Templates.Create(clientToken.Access, support.Template{
				Name:    "Gravity",
				Subject: "Coca cola {{.Subject}}",
				HTML:    "<h1>Rat</h1>{{.HTML}}<section>{{.Endorsement}}</section>",
				Text:    "Rat\n{{.Text}}\n{{.Endorsement}}",
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

		By("sending a notification to an organization", func() {
			status, responses, err := client.Notify.Organization(clientToken.Access, "org-123", support.Notify{
				KindID:  "organization-test",
				HTML:    "this is an organization test",
				Text:    "this is an organization test",
				Subject: "organization-subject",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(responses).To(HaveLen(3))

			for _, response := range responses {
				indexedResponses[response.Recipient] = response
			}

			response := indexedResponses["user-456"]
			Expect(response.Recipient).To(Equal("user-456"))
			Expect(response.Status).To(Equal("queued"))
			Expect(GUIDRegex.MatchString(response.NotificationID)).To(BeTrue())

			response = indexedResponses["user-789"]
			Expect(response.Recipient).To(Equal("user-789"))
			Expect(response.Status).To(Equal("queued"))
			Expect(GUIDRegex.MatchString(response.NotificationID)).To(BeTrue())

			response = indexedResponses["user-000"]
			Expect(response.Recipient).To(Equal("user-000"))
			Expect(response.Status).To(Equal("queued"))
			Expect(GUIDRegex.MatchString(response.NotificationID)).To(BeTrue())
		})

		By("confirming the messages were sent", func() {
			Eventually(func() int {
				return len(Servers.SMTP.Deliveries)
			}, 1*time.Second).Should(Equal(1))
			delivery := Servers.SMTP.Deliveries[0]

			env := application.NewEnvironment()
			Expect(delivery.Sender).To(Equal(env.Sender))
			Expect(delivery.Recipients).To(Equal([]string{"user-456@example.com"}))

			data := strings.Split(string(delivery.Data), "\n")
			Expect(data).To(ContainElement("X-CF-Client-ID: notifications-sender"))
			Expect(data).To(ContainElement("X-CF-Notification-ID: " + indexedResponses["user-456"].NotificationID))
			Expect(data).To(ContainElement("Subject: Coca cola organization-subject"))
			Expect(data).To(ContainElement("\t\t<h1>Rat</h1>this is an organization test<section>You received this message="))
			Expect(data).To(ContainElement(" because you belong to the \"notifications-service\" organization.</section>"))
			Expect(data).To(ContainElement("Rat"))
			Expect(data).To(ContainElement("this is an organization test"))
			Expect(data).To(ContainElement("You received this message because you belong to the \"notifications-service\" or="))
			Expect(data).To(ContainElement("ganization."))
		})
	})
})
