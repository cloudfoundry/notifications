package middleware_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/v1/web/middleware"
	"github.com/gorilla/mux"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rcrowley/go-metrics"
)

var _ = Describe("RequestCounter", func() {
	var (
		ware    middleware.RequestCounter
		request *http.Request
		writer  *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		var err error
		request, err = http.NewRequest("GET", "/clients/my-client/notifications/my-notification", nil)
		Expect(err).NotTo(HaveOccurred())

		writer = httptest.NewRecorder()
		matcher := mux.NewRouter()
		path := "/clients/{client_id}/notifications/{notification_id}"
		matcher.HandleFunc(path, func(http.ResponseWriter, *http.Request) {}).Name("GET " + path)

		ware = middleware.NewRequestCounter(matcher)
	})

	It("logs a request hit for a matching route", func() {
		result := ware.ServeHTTP(writer, request, nil)

		Expect(result).To(BeTrue())

		counter := metrics.GetOrRegisterCounter("notifications.web.GET./clients/:client_id/notifications/:notification_id", nil)

		Expect(counter.Count()).To(BeEquivalentTo(1))
	})
})
