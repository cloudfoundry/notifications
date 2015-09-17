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

		var results struct {
			Links struct {
				Self struct {
					Href string
				}

				Senders struct {
					Href string
				}

				Templates struct {
					Href string
				}
			} `json:"_links"`
		}

		status, err := client.DoTyped("GET", "/", nil, "", &results)
		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(http.StatusOK))

		Expect(results.Links.Self.Href).To(Equal("/"))
		Expect(results.Links.Senders.Href).To(Equal("/senders"))
		Expect(results.Links.Templates.Href).To(Equal("/templates"))
	})
})
