package acceptance

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/v2/acceptance/support"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("v2 API", func() {
	var (
		client *support.Client
	)

	BeforeEach(func() {
		client = support.NewClient(support.Config{
			Host:  Servers.Notifications.URL(),
			Trace: Trace,
		})
	})

	It("serves the correct API version number", func() {
		status, response, err := client.Do("GET", "/info", nil, "")
		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(http.StatusOK))
		Expect(response["version"]).To(Equal(float64(2)))
	})
})
