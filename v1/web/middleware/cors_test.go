package middleware_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/v1/web/middleware"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CORS", func() {
	var (
		writer     *httptest.ResponseRecorder
		request    *http.Request
		ware       middleware.CORS
		corsOrigin = "test-cors-origin"
	)

	Describe("ServeHTTP", func() {
		BeforeEach(func() {
			var err error

			writer = httptest.NewRecorder()
			request, err = http.NewRequest("OPTIONS", "/user_preferences", nil)
			if err != nil {
				panic(err)
			}

			ware = middleware.NewCORS(corsOrigin)
		})

		It("sets the correct CORS headers", func() {
			result := ware.ServeHTTP(writer, request, nil)

			Expect(result).To(BeTrue())
			Expect(writer.HeaderMap.Get("Access-Control-Allow-Origin")).To(Equal("test-cors-origin"))
			Expect(writer.HeaderMap.Get("Access-Control-Allow-Methods")).To(Equal("GET, PATCH"))
			Expect(writer.HeaderMap.Get("Access-Control-Allow-Headers")).To(Equal("Accept, Authorization, Content-Type"))
		})
	})
})
