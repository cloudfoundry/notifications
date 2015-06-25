package web_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/web"
	"github.com/gorilla/mux"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RouterPool", func() {
	var (
		pool    *web.RouterPool
		router1 *mux.Router
		router2 *mux.Router
	)

	BeforeEach(func() {
		pool = web.NewRouterPool()
		router1 = mux.NewRouter()
		router2 = mux.NewRouter()

		router1.HandleFunc("/router1", func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte("Hello from router1"))
		})
		router1.HandleFunc("/conflict", func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte("Conflict from router1"))
		})

		router2.HandleFunc("/router2", func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte("Hello from router2"))
		})
		router2.HandleFunc("/conflict", func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte("Conflict from router2"))
		})
	})

	It("serves traffic from multiple routers", func() {
		pool.Add(web.MuxToMatchableRouter(router1))
		pool.Add(web.MuxToMatchableRouter(router2))

		request, err := http.NewRequest("GET", "/router1", nil)
		Expect(err).NotTo(HaveOccurred())

		recorder := httptest.NewRecorder()
		pool.ServeHTTP(recorder, request)

		Expect(recorder.Body.String()).To(Equal("Hello from router1"))

		request, err = http.NewRequest("GET", "/router2", nil)
		Expect(err).NotTo(HaveOccurred())

		recorder = httptest.NewRecorder()
		pool.ServeHTTP(recorder, request)

		Expect(recorder.Body.String()).To(Equal("Hello from router2"))
	})

	It("follows order of inclusion for matching", func() {
		pool.Add(web.MuxToMatchableRouter(router1))
		pool.Add(web.MuxToMatchableRouter(router2))

		request, err := http.NewRequest("GET", "/conflict", nil)
		Expect(err).NotTo(HaveOccurred())

		recorder := httptest.NewRecorder()
		pool.ServeHTTP(recorder, request)

		Expect(recorder.Body.String()).To(Equal("Conflict from router1"))
	})

	Context("when there is no matching router", func() {
		It("handles the request with http.NotFound", func() {
			pool.Add(web.MuxToMatchableRouter(router1))
			pool.Add(web.MuxToMatchableRouter(router2))

			request, err := http.NewRequest("GET", "/missing", nil)
			Expect(err).NotTo(HaveOccurred())

			recorder := httptest.NewRecorder()
			pool.ServeHTTP(recorder, request)

			Expect(recorder.Code).To(Equal(http.StatusNotFound))
			Expect(recorder.Body.String()).To(Equal("404 page not found\n"))
		})
	})
})
