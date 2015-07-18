package v2

import (
	"github.com/cloudfoundry-incubator/notifications/acceptance/v2/support"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Notification types lifecycle", func() {
	var (
		client *support.Client
		token  uaa.Token
		sender support.Sender
	)

	BeforeEach(func() {
		client = support.NewClient(support.Config{
			Host: Servers.Notifications.URL(),
		})
		token = GetClientTokenFor("my-client", "uaa")
		var err error
		sender, err = client.Senders.Create("my-sender", token.Access)
		Expect(err).NotTo(HaveOccurred())
	})

	PIt("can create a new notification type", func() {
		By("creating a notification type", func() {
			notificationType, err := client.NotificationTypes.Create(sender.ID, "some-notification-type", "a great notification type", "", false, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(notificationType.Name).To(Equal("some-notification-type"))
			Expect(notificationType.Description).To(Equal("a great notification type"))
			Expect(notificationType.Critical).To(BeFalse())
			Expect(notificationType.TemplateID).To(BeEmpty())
		})
	})
})
