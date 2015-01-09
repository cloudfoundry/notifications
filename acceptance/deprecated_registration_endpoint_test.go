package acceptance

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/notifications/acceptance/support"
	"github.com/cloudfoundry-incubator/notifications/application"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("notifications can be registered, using the deprecated /registration endpoint", func() {
	It("registers a notification", func() {
		clientToken := GetClientTokenFor("notifications-sender")
		client := support.NewClient(Servers.Notifications)
		userID := "user-123"

		By("registering a notification with the deprecated endpoint", func() {
			body, err := json.Marshal(map[string]interface{}{
				"source_description": "Notifications Sender",
				"kinds": []map[string]string{
					{
						"id":          "acceptance-test",
						"description": "Acceptance Test",
					},
				},
			})
			if err != nil {
				panic(err)
			}

			request, err := http.NewRequest("PUT", Servers.Notifications.RegistrationPath(), bytes.NewBuffer(body))
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

		var response support.NotifyResponse

		By("sending a notifications to a user", func() {
			status, responses, err := client.Notify.User(clientToken.Access, userID, support.Notify{
				KindID:  "acceptance-test",
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
			}, 1*time.Second).Should(Equal(1))
			delivery := Servers.SMTP.Deliveries[0]

			env := application.NewEnvironment()
			Expect(delivery.Sender).To(Equal(env.Sender))
			Expect(delivery.Recipients).To(Equal([]string{"user-123@example.com"}))

			data := strings.Split(string(delivery.Data), "\n")
			Expect(data).To(ContainElement("X-CF-Client-ID: notifications-sender"))
			Expect(data).To(ContainElement("X-CF-Notification-ID: " + response.NotificationID))
			Expect(data).To(ContainElement("Subject: CF Notification: my-special-subject"))
			Expect(data).To(ContainElement(`        <p>This message was sent directly to you.</p><p>this is an acceptance%40test</p>`))
			Expect(data).To(ContainElement("hello from the acceptance test"))
		})
	})
})
