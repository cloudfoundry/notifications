package v1

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/v1/acceptance/support"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("v1 API", func() {
	var (
		client *support.Client
	)

	BeforeEach(func() {
		client = support.NewClient(Servers.Notifications.URL())
	})

	It("serves the correct API version number", func() {
		status, version, err := client.API.Version()
		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(http.StatusOK))
		Expect(version).To(Equal(1))
	})
})
