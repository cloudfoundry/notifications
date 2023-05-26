package v1

import (
	"net/http"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/cloudfoundry-incubator/notifications/v1/acceptance/support"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Preferences Endpoint", func() {
	var client *support.Client
	var clientToken uaa.Token
	var userToken uaa.Token
	var userGUID string
	var response support.NotifyResponse

	BeforeEach(func() {
		client = support.NewClient(Servers.Notifications.URL())
		clientToken = GetClientTokenFor("notifications-sender")
		userToken = GetUserTokenFor("user-123-code")
		userGUID = "user-123"

		By("registering a client with a notification", func() {
			status, err := client.Notifications.Register(clientToken.Access, support.RegisterClient{
				SourceName: "Notifications Sender",
				Notifications: map[string]support.RegisterNotification{
					"acceptance-test": {
						Description: "Acceptance Test",
					},
					"unsubscribe-acceptance-test": {
						Description: "Unsubscribe Acceptance Test",
					},
				},
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNoContent))
		})

		By("sending a notification to the user", func() {
			status, responses, err := client.Notify.User(clientToken.Access, userGUID, support.Notify{
				KindID:  "unsubscribe-acceptance-test",
				HTML:    "<p>this is an acceptance test</p>",
				Subject: "my-special-subject",
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(responses).To(HaveLen(1))

			response = responses[0]
			Expect(response.Status).To(Equal("queued"))
			Expect(response.Recipient).To(Equal(userGUID))
			Expect(GUIDRegex.MatchString(response.NotificationID)).To(BeTrue())
		})

		By("confirming that the messages were sent", func() {
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
			Expect(data).To(ContainElement("\t\t<p>This message was sent directly to you.</p><p>this is an acceptance test="))
			Expect(data).To(ContainElement("</p>"))
		})
	})

	It("user unsubscribes from a notification", func() {
		By("retrieving the current user preferences", func() {
			status, preferences, err := client.Preferences.Get(userToken.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(preferences.NotificationPreferences).To(HaveLen(2))

			Expect(preferences.NotificationPreferences).To(ContainElement(support.Preference{
				ClientID:                "notifications-sender",
				NotificationID:          "acceptance-test",
				Email:                   true,
				NotificationDescription: "Acceptance Test",
				SourceDescription:       "Notifications Sender",
			}))
			Expect(preferences.NotificationPreferences).To(ContainElement(support.Preference{
				ClientID:                "notifications-sender",
				NotificationID:          "unsubscribe-acceptance-test",
				Email:                   true,
				NotificationDescription: "Unsubscribe Acceptance Test",
				SourceDescription:       "Notifications Sender",
			}))
		})

		By("unsubscribing from a notification", func() {
			status, err := client.Preferences.Unsubscribe(userToken.Access, "notifications-sender", "unsubscribe-acceptance-test")
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNoContent))
		})

		By("confirming that the user is unsubscribed", func() {
			status, preferences, err := client.Preferences.Get(userToken.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(preferences.NotificationPreferences).To(HaveLen(2))

			Expect(preferences.NotificationPreferences).To(ContainElement(support.Preference{
				ClientID:                "notifications-sender",
				NotificationID:          "acceptance-test",
				Email:                   true,
				NotificationDescription: "Acceptance Test",
				SourceDescription:       "Notifications Sender",
			}))
			Expect(preferences.NotificationPreferences).To(ContainElement(support.Preference{
				ClientID:                "notifications-sender",
				NotificationID:          "unsubscribe-acceptance-test",
				Email:                   false,
				NotificationDescription: "Unsubscribe Acceptance Test",
				SourceDescription:       "Notifications Sender",
			}))
		})

		By("confirming that the user no longer receives messages for this notification", func() {
			Servers.SMTP.Reset()

			status, responses, err := client.Notify.User(clientToken.Access, userGUID, support.Notify{
				KindID:  "unsubscribe-acceptance-test",
				HTML:    "<p>this is an acceptance test</p>",
				Subject: "my-special-subject",
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(responses).To(HaveLen(1))

			response := responses[0]
			Expect(response.Status).To(Equal("queued"))
			Expect(response.Recipient).To(Equal(userGUID))
			Expect(GUIDRegex.MatchString(response.NotificationID)).To(BeTrue())
		})

		By("confirming that the email never gets delivered", func() {
			Consistently(func() int {
				return len(Servers.SMTP.Deliveries)
			}, 1*time.Second).Should(Equal(0))

			err := Servers.Notifications.WaitForJobsQueueToEmpty()
			Expect(err).NotTo(HaveOccurred())
		})

		By("re-subscribing to a notification", func() {
			status, err := client.Preferences.Subscribe(userToken.Access, "notifications-sender", "unsubscribe-acceptance-test")
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNoContent))
		})

		By("confirming that the user is subscribed", func() {
			status, preferences, err := client.Preferences.Get(userToken.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(preferences.NotificationPreferences).To(HaveLen(2))

			Expect(preferences.NotificationPreferences).To(ContainElement(support.Preference{
				ClientID:                "notifications-sender",
				NotificationID:          "acceptance-test",
				Email:                   true,
				NotificationDescription: "Acceptance Test",
				SourceDescription:       "Notifications Sender",
			}))
			Expect(preferences.NotificationPreferences).To(ContainElement(support.Preference{
				ClientID:                "notifications-sender",
				NotificationID:          "unsubscribe-acceptance-test",
				Email:                   true,
				NotificationDescription: "Unsubscribe Acceptance Test",
				SourceDescription:       "Notifications Sender",
			}))
		})

		By("confirming that the user receives messages for this notification", func() {
			Servers.SMTP.Reset()

			status, responses, err := client.Notify.User(clientToken.Access, userGUID, support.Notify{
				KindID:  "unsubscribe-acceptance-test",
				HTML:    "<p>this is an acceptance test</p>",
				Subject: "my-special-subject",
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(responses).To(HaveLen(1))

			response = responses[0]
			Expect(response.Status).To(Equal("queued"))
			Expect(response.Recipient).To(Equal(userGUID))
			Expect(GUIDRegex.MatchString(response.NotificationID)).To(BeTrue())
		})

		By("confirming that the email gets delivered", func() {
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
			Expect(data).To(ContainElement("\t\t<p>This message was sent directly to you.</p><p>this is an acceptance test="))
			Expect(data).To(ContainElement("</p>"))
		})
	})

	It("allows the user to globally unsubscribe", func() {
		By("retrieving the current user preferences", func() {
			status, preferences, err := client.Preferences.Get(userToken.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(preferences.NotificationPreferences).To(HaveLen(2))

			Expect(preferences.NotificationPreferences).To(ContainElement(support.Preference{
				ClientID:                "notifications-sender",
				NotificationID:          "acceptance-test",
				Email:                   true,
				NotificationDescription: "Acceptance Test",
				SourceDescription:       "Notifications Sender",
			}))
			Expect(preferences.NotificationPreferences).To(ContainElement(support.Preference{
				ClientID:                "notifications-sender",
				NotificationID:          "unsubscribe-acceptance-test",
				Email:                   true,
				NotificationDescription: "Unsubscribe Acceptance Test",
				SourceDescription:       "Notifications Sender",
			}))
		})

		By("globally unsubscribing from notifications", func() {
			status, err := client.Preferences.GlobalUnsubscribe(userToken.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNoContent))
		})

		By("confirming that the user is globally unsubscribed", func() {
			status, preferences, err := client.Preferences.Get(userToken.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(preferences.GlobalUnsubscribe).To(BeTrue())
		})

		By("confirming no longer receives notifications", func() {
			Servers.SMTP.Reset()

			status, responses, err := client.Notify.User(clientToken.Access, userGUID, support.Notify{
				KindID:  "acceptance-test",
				HTML:    "<p>this is an acceptance test</p>",
				Subject: "my-special-subject",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(responses).To(HaveLen(1))

			response := responses[0]
			Expect(response.Status).To(Equal("queued"))
			Expect(response.Recipient).To(Equal(userGUID))
			Expect(GUIDRegex.MatchString(response.NotificationID)).To(BeTrue())

			Consistently(func() int {
				return len(Servers.SMTP.Deliveries)
			}, 1*time.Second).Should(Equal(0))

			err = Servers.Notifications.WaitForJobsQueueToEmpty()
			Expect(err).NotTo(HaveOccurred())
		})

		By("re-subscribing globally to notifications", func() {
			status, err := client.Preferences.GlobalSubscribe(userToken.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNoContent))
		})

		By("confirming that the user is globally subscribed", func() {
			status, preferences, err := client.Preferences.Get(userToken.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(preferences.GlobalUnsubscribe).To(BeFalse())
		})

		By("confirming that the user now receives notifications", func() {
			Servers.SMTP.Reset()

			status, responses, err := client.Notify.User(clientToken.Access, userGUID, support.Notify{
				KindID:  "acceptance-test",
				HTML:    "<p>this is an acceptance test</p>",
				Subject: "my-special-subject",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(responses).To(HaveLen(1))

			response := responses[0]
			Expect(response.Status).To(Equal("queued"))
			Expect(response.Recipient).To(Equal(userGUID))
			Expect(GUIDRegex.MatchString(response.NotificationID)).To(BeTrue())

			Eventually(func() int {
				return len(Servers.SMTP.Deliveries)
			}, 10*time.Second).Should(Equal(1))
		})
	})
})
