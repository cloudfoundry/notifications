package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/cloudfoundry-incubator/notifications/v1/acceptance/support"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("notifications can be registered, using the deprecated /registration endpoint", func() {
	It("registers a notification", func() {
		var templateID string
		var response support.NotifyResponse

		clientID := "notifications-sender"
		notificationID := "acceptance-test"
		clientToken := GetClientTokenFor(clientID)
		client := support.NewClient(Servers.Notifications.URL())
		userID := "user-123"

		By("registering a notification with the deprecated endpoint", func() {
			body, err := json.Marshal(map[string]interface{}{
				"source_description": "Notifications Sender",
				"kinds": []map[string]string{
					{
						"id":          notificationID,
						"description": "Acceptance Test",
					},
				},
			})
			if err != nil {
				panic(err)
			}

			request, err := http.NewRequest("PUT", Servers.Notifications.URL()+"/registration", bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}

			request.Header.Set("Authorization", "Bearer "+clientToken.Access)

			response, err := http.DefaultClient.Do(request)
			if err != nil {
				panic(err)
			}

			// Confirm response status code looks ok
			Expect(response.StatusCode).To(Equal(http.StatusOK))
		})

		By("sending a notifications to a user", func() {
			status, responses, err := client.Notify.User(clientToken.Access, userID, support.Notify{
				KindID:  notificationID,
				HTML:    "<p>this is an acceptance%40test</p>",
				Text:    "hello from the acceptance test",
				Subject: "my-special-subject",
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))

			Expect(responses).To(HaveLen(1))
			response = responses[0]
			Expect(response.Status).To(Equal("queued"))
			Expect(response.Recipient).To(Equal(userID))
			Expect(GUIDRegex.MatchString(response.NotificationID)).To(BeTrue())
		})

		By("verifying that the message was sent", func() {
			Eventually(func() int {
				return len(Servers.SMTP.Deliveries)
			}, 10*time.Second).Should(Equal(1))
			delivery := Servers.SMTP.Deliveries[0]

			env, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())

			Expect(delivery.Sender).To(Equal(env.Sender))
			Expect(delivery.Recipients).To(Equal([]string{"user-123@example.com"}))

			data := strings.Split(string(delivery.Data), "\n")
			Expect(data).To(ContainElement("X-CF-Client-ID: notifications-sender"))
			Expect(data).To(ContainElement("X-CF-Notification-ID: " + response.NotificationID))
			Expect(data).To(ContainElement("Subject: CF Notification: my-special-subject"))
			Expect(data).To(ContainElement("\t\t<p>This message was sent directly to you.</p><p>this is an acceptance%40te="))
			Expect(data).To(ContainElement("st</p>"))
			Expect(data).To(ContainElement("hello from the acceptance test"))
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

		By("assigning the template to the client and notification", func() {
			status, err := client.Templates.AssignToClient(clientToken.Access, clientID, templateID)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNoContent))

			status, err = client.Templates.AssignToNotification(clientToken.Access, clientID, notificationID, templateID)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNoContent))
		})

		By("confirming that the client and notification have the assigned template", func() {
			status, notifications, err := client.Notifications.List(clientToken.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(notifications).To(HaveLen(1))
			Expect(notifications[clientID].Template).To(Equal(templateID))

			Expect(notifications[clientID].Notifications).To(HaveLen(1))
			Expect(notifications[clientID].Notifications[notificationID].Template).To(Equal(templateID))
		})

		By("re-registering the client with the deprecated endpoint", func() {
			body, err := json.Marshal(map[string]interface{}{
				"source_description": "Notifications Sender",
				"kinds": []map[string]string{
					{
						"id":          notificationID,
						"description": "Acceptance Test",
					},
				},
			})
			if err != nil {
				panic(err)
			}

			request, err := http.NewRequest("PUT", Servers.Notifications.URL()+"/registration", bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}

			request.Header.Set("Authorization", "Bearer "+clientToken.Access)

			response, err := http.DefaultClient.Do(request)
			if err != nil {
				panic(err)
			}

			// Confirm response status code looks ok
			Expect(response.StatusCode).To(Equal(http.StatusOK))
		})

		By("confirming that the client and notification continue to have the assigned template", func() {
			status, notifications, err := client.Notifications.List(clientToken.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(notifications).To(HaveLen(1))
			Expect(notifications[clientID].Template).To(Equal(templateID))

			Expect(notifications[clientID].Notifications).To(HaveLen(1))
			Expect(notifications[clientID].Notifications[notificationID].Template).To(Equal(templateID))
		})
	})
})
