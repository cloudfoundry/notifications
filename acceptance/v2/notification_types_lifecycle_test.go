package v2

import (
	"fmt"
	"net/http"

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

	It("can create and show a new notification type", func() {
		var notificationType support.NotificationType
		var err error
		By("creating a notification type", func() {
			notificationType, err = client.NotificationTypes.Create(sender.ID, "some-notification-type", "a great notification type", "", false, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(notificationType.Name).To(Equal("some-notification-type"))
			Expect(notificationType.Description).To(Equal("a great notification type"))
			Expect(notificationType.Critical).To(BeFalse())
			Expect(notificationType.TemplateID).To(BeEmpty())
		})

		By("creating it again with the same name", func() {
			notificationType, err = client.NotificationTypes.Create(sender.ID, "some-notification-type", "another great notification type", "", false, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(notificationType.Name).To(Equal("some-notification-type"))
			Expect(notificationType.Description).To(Equal("a great notification type"))
		})

		By("showing the newly created notification type", func() {
			gottenNotificationType, err := client.NotificationTypes.Show(sender.ID, notificationType.ID, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(gottenNotificationType.Name).To(Equal("some-notification-type"))
			Expect(gottenNotificationType.Description).To(Equal("a great notification type"))
		})
	})

	It("does not leak the existence of unauthorized notification types", func() {
		var notificationType support.NotificationType
		var otherSender support.Sender
		var err error

		By("creating a notification type belonging to 'my-sender'", func() {
			notificationType, err = client.NotificationTypes.Create(sender.ID, "some-notification-type", "a great notification type", "", false, token.Access)
			Expect(err).NotTo(HaveOccurred())
		})

		By("creating a sender that is not 'my-sender'", func() {
			otherSender, err = client.Senders.Create("some-other-sender", token.Access)
			Expect(err).NotTo(HaveOccurred())
		})

		By("verifying that you cannot get a notification type belonging to a different sender", func() {
			_, err := client.NotificationTypes.Show(otherSender.ID, notificationType.ID, token.Access)

			Expect(err.(support.NotFoundError).Status).To(Equal(http.StatusNotFound))
			expectedErrorMessage := fmt.Sprintf("{\"error\": \"notification type %s not found\"}", notificationType.ID)
			Expect(err.(support.NotFoundError).Body).To(MatchJSON(expectedErrorMessage))
		})
	})
})
