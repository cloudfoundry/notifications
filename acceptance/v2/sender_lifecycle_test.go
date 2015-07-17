package v2

import (
	"github.com/cloudfoundry-incubator/notifications/acceptance/v2/support"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sender lifecycle", func() {
	var (
		client *support.Client
		token  uaa.Token
	)

	BeforeEach(func() {
		client = support.NewClient(support.Config{
			Host: Servers.Notifications.URL(),
		})
		token = GetClientTokenFor("my-client", "uaa")
	})

	It("can create and read a new sender", func() {
		var sender support.Sender

		By("creating a sender", func() {
			var err error
			sender, err = client.Senders.Create("My Cool App", token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(sender.Name).To(Equal("My Cool App"))
		})

		By("getting the sender", func() {
			retrieved_sender, err := client.Senders.Get(sender.ID, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(retrieved_sender.Name).To(Equal("My Cool App"))
			Expect(retrieved_sender.ID).To(Equal(sender.ID))
		})
	})
})
