package middleware_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/v2/web/middleware"
	"github.com/gorilla/mux"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RequestCounter", func() {
	var (
		ware    middleware.RequestCounter
		request *http.Request
		writer  *httptest.ResponseRecorder
		buffer  *bytes.Buffer
	)

	BeforeEach(func() {
		var err error
		request, err = http.NewRequest("GET", "/clients/my-client/notifications/my-notification", nil)
		Expect(err).NotTo(HaveOccurred())

		writer = httptest.NewRecorder()
		matcher := mux.NewRouter()
		path := "/clients/{client_id}/notifications/{notification_id}"
		matcher.HandleFunc(path, func(http.ResponseWriter, *http.Request) {}).Name("GET " + path)
		buffer = bytes.NewBuffer([]byte{})

		ware = middleware.NewRequestCounter(matcher, log.New(buffer, "", 0))
	})

	It("logs a request hit for a matching route", func() {
		result := ware.ServeHTTP(writer, request, nil)

		Expect(result).To(BeTrue())

		Expect(buffer).To(MatchJSON(`{
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
