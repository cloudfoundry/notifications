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

var _ = Describe("Sending notifications to all users in a space", func() {
	It("sends a notification to each user in a space", func() {
		var templateID string
		client := support.NewClient(Servers.Notifications.URL())
		clientID := "notifications-sender"
		clientToken := GetClientTokenFor(clientID)
		spaceID := "space-123"
		indexedResponses := map[string]support.NotifyResponse{}

		By("registering a client with a notification", func() {
			status, err := client.Notifications.Register(clientToken.Access, support.RegisterClient{
				SourceName: "Notifications Sender",
				Notifications: map[string]support.RegisterNotification{
					"space-test": {
						Description: "Space Test",
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
				Name:    "Men in Black",
				Subject: "Aliens {{.Subject}}",
				HTML:    "<h1>Dogs</h1>{{.HTML}}<h2>{{.Endorsement}}</h2>",
				Text:    "Dogs\n{{.Text}}\n{{.Endorsement}}",
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

		By("sending a notification to the users of a space", func() {
			status, responses, err := client.Notify.Space(clientToken.Access, spaceID, support.Notify{
				KindID:  "space-test",
				HTML:    "this is a space test",
				Text:    "this is a space test",
				Subject: "space-subject",
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
			Expect(response.VCAPRequestID).To(Equal("some-totally-fake-vcap-request-id"))

			response = indexedResponses["user-789"]
			Expect(response.Recipient).To(Equal("user-789"))
			Expect(response.Status).To(Equal("queued"))
			Expect(GUIDRegex.MatchString(response.NotificationID)).To(BeTrue())
			Expect(response.VCAPRequestID).To(Equal("some-totally-fake-vcap-request-id"))

			response = indexedResponses["user-000"]
			Expect(response.Recipient).To(Equal("user-000"))
			Expect(response.Status).To(Equal("queued"))
			Expect(GUIDRegex.MatchString(response.NotificationID)).To(BeTrue())
			Expect(response.VCAPRequestID).To(Equal("some-totally-fake-vcap-request-id"))
		})

		By("confirming the messages were sent", func() {
			Eventually(func() int {
				return len(Servers.SMTP.Deliveries)
			}, 10*time.Second).Should(Equal(1))
			delivery := Servers.SMTP.Deliveries[0]

			env, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())

			Expect(delivery.Sender).To(Equal(env.Sender))
			Expect(delivery.Recipients).To(Equal([]string{"user-456@example.com"}))

			data := strings.Split(string(delivery.Data), "\n")
			Expect(data).To(ContainElement("X-CF-Client-ID: notifications-sender"))
			Expect(data).To(ContainElement("X-CF-Notification-ID: " + indexedResponses["user-456"].NotificationID))
			Expect(data).To(ContainElement("Subject: Aliens space-subject"))
			Expect(data).To(ContainElement("\t\t<h1>Dogs</h1>this is a space test<h2>You received this message because you="))
			Expect(data).To(ContainElement(` belong to the &#34;notifications-service&#34; space in the &#34;notificatio=`))
			Expect(data).To(ContainElement(`ns-service&#34; organization.</h2>`))
			Expect(data).To(ContainElement("Dogs"))
			Expect(data).To(ContainElement("this is a space test"))
			Expect(data).To(ContainElement(`You received this message because you belong to the "notifications-service" =`))
			Expect(data).To(ContainElement(`space in the "notifications-service" organization.`))
		})
	})

	It("returns a 404 if the space cannot be found", func() {
		client := support.NewClient(Servers.Notifications.URL())
		clientID := "notifications-sender"
		clientToken := GetClientTokenFor(clientID)
		spaceID := "banana"

		status, _, err := client.Notify.Space(clientToken.Access, spaceID, support.Notify{
			KindID:  "space-test",
			HTML:    "this is a space test",
			Text:    "this is a space test",
			Subject: "space-subject",
		})

		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(http.StatusNotFound))
	})
})
