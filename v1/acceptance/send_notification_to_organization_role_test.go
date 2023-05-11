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

var _ = Describe("Sending notifications to users with certain roles in an organization", func() {
	var (
		templateID  string
		clientID    string
		clientToken uaa.Token
		client      *support.Client
	)

	BeforeEach(func() {
		clientID = "notifications-sender"
		clientToken = GetClientTokenFor(clientID)
		client = support.NewClient(Servers.Notifications.URL())
		Servers.SMTP.Reset()

		By("registering a notification", func() {
			status, err := client.Notifications.Register(clientToken.Access, support.RegisterClient{
				SourceName: "Notifications Sender",
				Notifications: map[string]support.RegisterNotification{
					"organization-role-test": {
						Description: "Organization Role Test",
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
				Name:    "ET",
				Subject: "Phone home {{.Subject}}",
				HTML:    "<h1>Cat</h1>{{.HTML}}<header>{{.Endorsement}}</header>",
				Text:    "Cat\n{{.Text}}\n{{.Endorsement}}",
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
	})

	It("sends a notification to each OrgManager in an organization", func() {
		var response support.NotifyResponse

		By("sending a notification to the OrgManager role", func() {
			status, responses, err := client.Notify.OrganizationRole(clientToken.Access, "org-123", "OrgManager", support.Notify{
				KindID:  "organization-role-test",
				HTML:    "this is another organization role test",
				Text:    "this is an organization role test",
				Subject: "organization-role-subject",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(responses).To(HaveLen(1))

			response = responses[0]
			Expect(response.Recipient).To(Equal("user-456"))
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
			Expect(data).To(ContainElement("X-CF-Notification-ID: " + response.NotificationID))
			Expect(data).To(ContainElement("Subject: Phone home organization-role-subject"))
			Expect(data).To(ContainElement("Cat"))
			Expect(data).To(ContainElement("this is an organization role test"))
			Expect(data).To(ContainElement(`You received this message because you are an OrgManager in the "notification=`))
			Expect(data).To(ContainElement(`s-service" organization.`))
			Expect(data).To(ContainElement("\t\t<h1>Cat</h1>this is another organization role test<header>You received thi="))
			Expect(data).To(ContainElement(`s message because you are an OrgManager in the &#34;notifications-service&#3=`))
			Expect(data).To(ContainElement(`4; organization.</header>`))
		})
	})

	It("sends a notification to each auditor in an organization", func() {
		var response support.NotifyResponse

		By("sending a notification to the OrgAuditor role", func() {
			status, responses, err := client.Notify.OrganizationRole(clientToken.Access, "org-123", "OrgAuditor", support.Notify{
				KindID:  "organization-role-test",
				HTML:    "this is another organization role test",
				Text:    "this is an organization role test",
				Subject: "organization-role-subject",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(responses).To(HaveLen(1))

			response = responses[0]

			Expect(response.Recipient).To(Equal("user-123"))
			Expect(response.Status).To(Equal("queued"))
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
			Expect(data).To(ContainElement("Subject: Phone home organization-role-subject"))
			Expect(data).To(ContainElement("Cat"))
			Expect(data).To(ContainElement("this is an organization role test"))
			Expect(data).To(ContainElement(`You received this message because you are an OrgAuditor in the "notification=`))
			Expect(data).To(ContainElement(`s-service" organization.`))
			Expect(data).To(ContainElement("\t\t<h1>Cat</h1>this is another organization role test<header>You received thi="))
			Expect(data).To(ContainElement("s message because you are an OrgAuditor in the &#34;notifications-service&#3="))
			Expect(data).To(ContainElement("4; organization.</header>"))
		})
	})

	It("sends a notification to each billing manager in an organization", func() {
		var response support.NotifyResponse

		By("sending a notification to the BillingManager role", func() {
			status, responses, err := client.Notify.OrganizationRole(clientToken.Access, "org-123", "BillingManager", support.Notify{
				KindID:  "organization-role-test",
				HTML:    "this is another organization role test",
				Text:    "this is an organization role test",
				Subject: "organization-role-subject",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(responses).To(HaveLen(1))

			response = responses[0]

			Expect(response.Recipient).To(Equal("user-111"))
			Expect(response.Status).To(Equal("queued"))
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
			Expect(delivery.Recipients).To(Equal([]string{"user-111@example.com"}))

			data := strings.Split(string(delivery.Data), "\n")
			Expect(data).To(ContainElement("X-CF-Client-ID: notifications-sender"))
			Expect(data).To(ContainElement("X-CF-Notification-ID: " + response.NotificationID))
			Expect(data).To(ContainElement("Subject: Phone home organization-role-subject"))
			Expect(data).To(ContainElement("Cat"))
			Expect(data).To(ContainElement("this is an organization role test"))
			Expect(data).To(ContainElement(`You received this message because you are an BillingManager in the "notifica=`))
			Expect(data).To(ContainElement(`tions-service" organization.`))
			Expect(data).To(ContainElement("\t\t<h1>Cat</h1>this is another organization role test<header>You received thi="))
			Expect(data).To(ContainElement("s message because you are an BillingManager in the &#34;notifications-servic="))
			Expect(data).To(ContainElement("e&#34; organization.</header>"))
		})
	})

	It("sends a notification to an invalid role in an organization", func() {
		By("sending a notification to an invalid role", func() {
			status, _, err := client.Notify.OrganizationRole(clientToken.Access, "org-123", "bad-role", support.Notify{
				KindID:  "organization-role-test",
				HTML:    "this is another organization role test",
				Text:    "this is an organization role test",
				Subject: "organization-role-subject",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(422))
		})
	})
})
