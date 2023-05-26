package v1

import (
	"encoding/json"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("metrics endpoint", func() {
	It("returns 200 and exposes metrics", func() {
		resp, err := http.Get(Servers.Notifications.URL() + "/debug/metrics")

		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		var body map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&body)
		Expect(err).ToNot(HaveOccurred())
		defer resp.Body.Close()
		Expect(body["notifications.web.GET./info"]).To(BeNumerically( ">=", 1))
	})
})
