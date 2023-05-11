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

var _ = Describe("Sending notifications to users with certain scopes", func() {
	It("sends a notification to each user with the scope", func() {
		var templateID string
		indexedResponses := map[string]support.NotifyResponse{}

		client := support.NewClient(Servers.Notifications.URL())
		clientID := "notifications-sender"
		clientToken := GetClientTokenFor(clientID)
		scope := "this.scope"

		By("registering a client with a notification", func() {
			status, err := client.Notifications.Register(clientToken.Access, support.RegisterClient{
				SourceName: "Notifications Sender",
				Notifications: map[string]support.RegisterNotification{
					"scope-test": {
						Description: "Scope Test",
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
				Name:    "Frozen",
				Subject: "Food {{.Subject}}",
				HTML:    "<h1>Fish</h1>{{.HTML}}<b>{{.Endorsement}}</b>",
				Text:    "Fish\n{{.Text}}\n{{.Endorsement}}",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusCreated))
			Expect(templateID).NotTo(Equal(""))
		})

		By("assigning the template to the client", func() {
			status, err := client.Templates.AssignToClient(clientToken.Access, clientID, templateID)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNoContent))
		})

		By("sending a notification to all users with a UAA scope", func() {
			status, responses, err := client.Notify.Scope(clientToken.Access, scope, support.Notify{
				KindID:  "scope-test",
				HTML:    "this is a scope test",
				Text:    "this is a scope test",
				Subject: "scope-subject",
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(responses).To(HaveLen(1))

			for _, response := range responses {
				indexedResponses[response.Recipient] = response
			}
		})

		By("confirming that the messages were delivered", func() {
			response := indexedResponses["user-369"]
			Expect(response.Recipient).To(Equal("user-369"))
			Expect(response.Status).To(Equal("queued"))
			Expect(GUIDRegex.MatchString(response.NotificationID)).To(BeTrue())
			Expect(response.VCAPRequestID).To(Equal("some-totally-fake-vcap-request-id"))

			Eventually(func() int {
				return len(Servers.SMTP.Deliveries)
			}, 10*time.Second).Should(Equal(1))
			delivery := Servers.SMTP.Deliveries[0]

			env, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())

			Expect(delivery.Sender).To(Equal(env.Sender))
			Expect(delivery.Recipients).To(Equal([]string{"user-369@example.com"}))

			data := strings.Split(string(delivery.Data), "\n")
			Expect(data).To(ContainElement("X-CF-Client-ID: notifications-sender"))
			Expect(data).To(ContainElement("X-CF-Notification-ID: " + indexedResponses["user-369"].NotificationID))
			Expect(data).To(ContainElement("Subject: Food scope-subject"))
			Expect(data).To(ContainElement("\t\t<h1>Fish</h1>this is a scope test<b>You received this message because you ="))
			Expect(data).To(ContainElement("have the this.scope scope.</b>"))
			Expect(data).To(ContainElement("this is a scope test"))
			Expect(data).To(ContainElement("Fish"))
			Expect(data).To(ContainElement("You received this message because you have the this.scope scope."))
		})
	})

	It("does not send messages to a user with a default scope", func() {
		var templateID string
		client := support.NewClient(Servers.Notifications.URL())
		clientID := "notifications-sender"
		clientToken := GetClientTokenFor(clientID)
		scope := "openid"

		By("registering a client with a notification", func() {
			status, err := client.Notifications.Register(clientToken.Access, support.RegisterClient{
				SourceName: "Notifications Sender",
				Notifications: map[string]support.RegisterNotification{
					"scope-test": {
						Description: "Scope Test",
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
				Name:    "Frozen",
				Subject: "Food {{.Subject}}",
				HTML:    "<h1>Fish</h1>{{.HTML}}<b>{{.Endorsement}}</b>",
				Text:    "Fish\n{{.Text}}\n{{.Endorsement}}",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusCreated))
			Expect(templateID).NotTo(Equal(""))
		})

		By("assigning the template to the client", func() {
			status, err := client.Templates.AssignToClient(clientToken.Access, clientID, templateID)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNoContent))
		})

		By("sending a notification to all users with a UAA scope", func() {
			status, _, err := client.Notify.Scope(clientToken.Access, scope, support.Notify{
				KindID:  "scope-test",
				HTML:    "this is a scope test",
				Text:    "this is a scope test",
				Subject: "scope-subject",
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNotAcceptable))
		})
	})
})
