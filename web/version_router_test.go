package web_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/web"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("VersionRouter", func() {
	var (
		v1Called bool
		v3Called bool
		router   web.VersionRouter
		request  *http.Request
		writer   *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		var err error

		v1Called = false
		v3Called = false
		writer = httptest.NewRecorder()
		request, err = http.NewRequest("GET", "/", nil)
		Expect(err).NotTo(HaveOccurred())
		router = web.VersionRouter{
			1: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				v1Called = true
			}),
			3: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				v3Called = true
			}),
		}
	})

	It("uses the X-NOTIFICATIONS-VERSION header to choose a version for the API", func() {
		request.Header.Set("X-NOTIFICATIONS-VERSION", "3")
		router.ServeHTTP(writer, request)

		Expect(v3Called).To(BeTrue())
		Expect(v1Called).To(BeFalse())
	})

	It("it defaults to v1 if the X-NOTIFICATIONS-VERSION header is absent", func() {
		router.ServeHTTP(writer, request)

		Expect(v1Called).To(BeTrue())
		Expect(v3Called).To(BeFalse())
	})

	It("returns a 404 if the version number is not an integer", func() {
		request.Header.Set("X-NOTIFICATIONS-VERSION", "banana")
		router.ServeHTTP(writer, request)

		Expect(v1Called).To(BeFalse())
		Expect(v3Called).To(BeFalse())
		Expect(writer.Code).To(Equal(http.StatusNotFound))
	})

	It("returns a 404 if the version number is not a version we support", func() {
		request.Header.Set("X-NOTIFICATIONS-VERSION", "42")
		router.ServeHTTP(writer, request)

		Expect(v1Called).To(BeFalse())
		Expect(v3Called).To(BeFalse())
		Expect(writer.Code).To(Equal(http.StatusNotFound))
	})
})
