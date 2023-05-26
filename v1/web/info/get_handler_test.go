package info_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/v1/web/info"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetHandler", func() {
	Describe("ServeHTTP", func() {
		var handler info.GetHandler

		BeforeEach(func() {
			handler = info.NewGetHandler()
		})

		It("returns a 200 response code and an empty JSON body", func() {
			writer := httptest.NewRecorder()
			request, err := http.NewRequest("GET", "/info", nil)
			if err != nil {
				panic(err)
			}

			handler.ServeHTTP(writer, request, nil)

			Expect(writer.Code).To(Equal(http.StatusOK))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"version": 1
			}`))
		})
	})
})
