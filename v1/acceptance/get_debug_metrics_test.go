package v1

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("/debug/metrics endpoint exposes metrics", func() {
	It("returns a 200", func() {
		resp, err := http.Get(Servers.Notifications.URL() + "/debug/metrics")

		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
	})
})
