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

var _ = Describe("Default Template", func() {
	var (
		client      *support.Client
		clientToken uaa.Token
	)

	BeforeEach(func() {
		clientToken = GetClientTokenFor("notifications-admin")
		client = support.NewClient(Servers.Notifications.URL())
	})

	It("can retrieve the default template", func() {
		status, template, err := client.Templates.Default.Get(clientToken.Access)
		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(http.StatusOK))
		Expect(template).To(Equal(support.Template{
			Name:     "Default Template",
			Subject:  "CF Notification: {{.Subject}}",
			HTML:     "<p>{{.Endorsement}}</p>{{.HTML}}",
			Text:     "{{.Endorsement}}\n{{.Text}}",
			Metadata: map[string]interface{}{},
		}))
	})

	It("gracefully handles missing metadata for a default-template update", func() {
		By("setting the default template without a metadata field", func() {
			status, err := client.Templates.Default.Update(clientToken.Access, support.Template{
				Name:    "Fail Hard - The Sequel",
				Subject: "Failed: {{.Subject}}",
				HTML:    "<h1>THIS IS FAILUREEE!!!</h1>",
				Text:    "Y U NO Fail!!!",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNoContent))
		})

		By("verifying that metadata is defaulted to {}", func() {
			status, template, err := client.Templates.Default.Get(clientToken.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(template).To(Equal(support.Template{
				Name:     "Fail Hard - The Sequel",
				Subject:  "Failed: {{.Subject}}",
				HTML:     "<h1>THIS IS FAILUREEE!!!</h1>",
				Text:     "Y U NO Fail!!!",
				Metadata: map[string]interface{}{},
			}))
		})
	})

	It("can edit the default template", func() {
		By("editing the default template", func() {
			status, err := client.Templates.Default.Update(clientToken.Access, support.Template{
				Name:    "A Whole New Template",
				Subject: "Updated: {{.Subject}}",
				HTML:    "<h1>Updated!!!</h1>",
				Text:    "Updated!!!",
				Metadata: map[string]interface{}{
					"smurf": "favorite",
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNoContent))
		})

		By("verifying that the default template was updated", func() {
			status, template, err := client.Templates.Default.Get(clientToken.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(template).To(Equal(support.Template{
				Name:    "A Whole New Template",
				Subject: "Updated: {{.Subject}}",
				HTML:    "<h1>Updated!!!</h1>",
				Text:    "Updated!!!",
				Metadata: map[string]interface{}{
					"smurf": "favorite",
				},
			}))
		})

		By("restarting the notifications service", func() {
			Servers.Notifications.Restart()
			client = support.NewClient(Servers.Notifications.URL())
		})

		By("verifying that the default template still displays the overridden values", func() {
			status, template, err := client.Templates.Default.Get(clientToken.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(template).To(Equal(support.Template{
				Name:    "A Whole New Template",
				Subject: "Updated: {{.Subject}}",
				HTML:    "<h1>Updated!!!</h1>",
				Text:    "Updated!!!",
				Metadata: map[string]interface{}{
					"smurf": "favorite",
				},
			}))
		})
	})

	It("can send a notification with the default template", func() {
		var response support.NotifyResponse

		By("sending a notification to an email address", func() {
			status, responses, err := client.Notify.Email(clientToken.Access, "John User <user@example.com>", support.Notify{
				KindID:  "acceptance-test",
				HTML:    "<header>this is an acceptance test</header>",
				Subject: "my-special-subject",
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
			}, 10*time.Second).Should(Equal(1))
			delivery := Servers.SMTP.Deliveries[0]

			env, err := application.NewEnvironment()
			Expect(err).NotTo(HaveOccurred())

			Expect(delivery.Sender).To(Equal(env.Sender))
			Expect(delivery.Recipients).To(Equal([]string{"user@example.com"}))

			data := strings.Split(string(delivery.Data), "\n")
			Expect(data).To(ContainElement("X-CF-Client-ID: notifications-admin"))
			Expect(data).To(ContainElement("X-CF-Notification-ID: " + response.NotificationID))
			Expect(data).To(ContainElement("Subject: CF Notification: my-special-subject"))
			Expect(data).To(ContainElement("\t\t<p>This message was sent directly to your email address.</p><header>this i="))
			Expect(data).To(ContainElement("s an acceptance test</header>"))
		})
	})
})
