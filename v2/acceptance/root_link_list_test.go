package acceptance

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/v2/acceptance/support"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("v2 Root Link List", func() {
	var (
		client *support.Client
	)

	BeforeEach(func() {
		client = support.NewClient(support.Config{
			Host:  Servers.Notifications.URL(),
			Trace: Trace,
		})
	})

	It("returns the list of root links", func() {
		status, response, err := client.Do("GET", "/", nil, "")
		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(http.StatusOK))

		Expect(response["_links"]).To(HaveKeyWithValue("self", map[string]interface{}{
			"href": "/",
		}))
		Expect(response["_links"]).To(HaveKeyWithValue("senders", map[string]interface{}{
			"href": "/senders",
		}))
		Expect(response["_links"]).To(HaveKeyWithValue("templates", map[string]interface{}{
			"href": "/templates",
		}))
	})
})
