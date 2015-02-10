package middleware_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/gorilla/mux"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RequestCounter", func() {
	var ware middleware.RequestCounter
	var request *http.Request
	var writer *httptest.ResponseRecorder
	var metricsLogger *log.Logger
	var buffer *bytes.Buffer

	BeforeEach(func() {
		var err error
		metricsLogger = metrics.Logger
		request, err = http.NewRequest("GET", "/clients/my-client/notifications/my-notification", nil)
		if err != nil {
			panic(err)
		}
		writer = httptest.NewRecorder()
		buffer = bytes.NewBuffer([]byte{})
		metrics.Logger = log.New(buffer, "", 0)
		matcher := mux.NewRouter()
		path := "/clients/{client_id}/notifications/{notification_id}"
		matcher.HandleFunc(path, func(http.ResponseWriter, *http.Request) {}).Name("GET " + path)

		ware = middleware.NewRequestCounter(matcher)
	})

	AfterEach(func() {
		metrics.Logger = metricsLogger
	})

	It("logs a request hit for a matching route", func() {
		result := ware.ServeHTTP(writer, request, nil)

		Expect(result).To(BeTrue())

		metric := strings.TrimPrefix(buffer.String(), "[METRIC] ")
		Expect(metric).To(MatchJSON(`{
			"kind": "counter",
			"payload": {
				"name": "notifications.web",
				"tags": {
					"endpoint": "GET/clients/:client_id/notifications/:notification_id"
				}
			}
		}`))
	})
})
