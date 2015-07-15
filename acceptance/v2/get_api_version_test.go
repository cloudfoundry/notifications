package v2

import (
	"github.com/cloudfoundry-incubator/notifications/acceptance/v2/support"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("v2 API", func() {
	var (
		client *support.Client
	)

	BeforeEach(func() {
		client = support.NewClient(support.Config{
			Host: Servers.Notifications.URL(),
		})
	})

	It("serves the correct API version number", func() {
		version, err := client.API.Version()
		Expect(err).NotTo(HaveOccurred())
		Expect(version).To(Equal(2))
	})
})
