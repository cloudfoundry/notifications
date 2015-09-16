package root_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/v2/web/root"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetHandler", func() {
	Describe("ServeHTTP", func() {
		var handler root.GetHandler

		BeforeEach(func() {
			handler = root.NewGetHandler()
		})

		It("returns a 200 response code and a JSON body describing the resources", func() {
			writer := httptest.NewRecorder()
			request, err := http.NewRequest("GET", "/", nil)
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, nil)

			Expect(writer.Code).To(Equal(http.StatusOK))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"_links": {
					"self": {
						"href": "/"
					},
					"senders": {
						"href": "/senders"
					},
					"templates": {
						"href": "/templates"
					}
				}
			}`))
		})
	})
})
