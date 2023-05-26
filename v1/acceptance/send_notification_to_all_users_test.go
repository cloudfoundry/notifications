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

var _ = Describe("Send a notification to all users of UAA", func() {
	It("sends an email notification to all users of UAA", func() {
		var templateID string
		indexedResponses := map[string]support.NotifyResponse{}
		clientID := "notifications-sender"
		clientToken := GetClientTokenFor(clientID)
		client := support.NewClient(Servers.Notifications.URL())

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

		By("creating a template", func() {
			var status int
			var err error
			status, templateID, err = client.Templates.Create(clientToken.Access, support.Template{
				Name:    "Jurassic Park",
				Subject: "Genetics {{.Subject}}",
				HTML:    "<h1>T-Rex</h1>{{.HTML}}<b>{{.Endorsement}}</b>",
				Text:    "T-Rex\n{{.Text}}\n{{.Endorsement}}",
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

		By("sending a notification to all users", func() {
			status, responses, err := client.Notify.AllUsers(clientToken.Access, support.Notify{
				KindID:  "acceptance-test",
				HTML:    "<p>this is an acceptance-test</p>",
				Text:    "oh no!",
				Subject: "gone awry",
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))

			Expect(responses).To(HaveLen(2))

			for _, resp := range responses {
				indexedResponses[resp.Recipient] = resp
			}

			responseItem := indexedResponses["091b6583-0933-4d17-a5b6-66e54666c88e"]
			Expect(responseItem.Recipient).To(Equal("091b6583-0933-4d17-a5b6-66e54666c88e"))
			Expect(responseItem.Status).To(Equal("queued"))
			Expect(GUIDRegex.MatchString(responseItem.NotificationID)).To(BeTrue())
			Expect(responseItem.VCAPRequestID).To(Equal("some-totally-fake-vcap-request-id"))

			responseItem = indexedResponses["943e6076-b1a5-4404-811b-a1ee9253bf56"]
			Expect(responseItem.Recipient).To(Equal("943e6076-b1a5-4404-811b-a1ee9253bf56"))
			Expect(responseItem.Status).To(Equal("queued"))
			Expect(GUIDRegex.MatchString(responseItem.NotificationID)).To(BeTrue())
			Expect(responseItem.VCAPRequestID).To(Equal("some-totally-fake-vcap-request-id"))
		})

		By("confirming the messages were sent", func() {
			Eventually(func() int {
				return len(Servers.SMTP.Deliveries)
			}, 10*time.Second).Should(Equal(2))

			recipients := []string{Servers.SMTP.Deliveries[0].Recipients[0], Servers.SMTP.Deliveries[1].Recipients[0]}
			Expect(recipients).To(ConsistOf([]string{"why-email@example.com", "slayer@example.com"}))

			var recipientIndex int
			if Servers.SMTP.Deliveries[0].Recipients[0] == "why-email@example.com" {
				recipientIndex = 0
			} else {
				recipientIndex = 1
			}

			env, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())

			delivery := Servers.SMTP.Deliveries[recipientIndex]
			Expect(delivery.Sender).To(Equal(env.Sender))

			data := strings.Split(string(delivery.Data), "\n")
			Expect(data).To(ContainElement("X-CF-Client-ID: notifications-sender"))
			Expect(data).To(ContainElement("X-CF-Notification-ID: " + indexedResponses["091b6583-0933-4d17-a5b6-66e54666c88e"].NotificationID))
			Expect(data).To(ContainElement("Subject: Genetics gone awry"))
			Expect(data).To(ContainElement("\t\t<h1>T-Rex</h1><p>this is an acceptance-test</p><b>This message was sent to="))
			Expect(data).To(ContainElement(" everyone.</b>"))
			Expect(data).To(ContainElement("T-Rex"))
			Expect(data).To(ContainElement("oh no!"))
			Expect(data).To(ContainElement("This message was sent to everyone."))
		})
	})
})
