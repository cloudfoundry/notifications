package middleware_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RequestLogging", func() {
	var (
		ware      middleware.RequestLogging
		request   *http.Request
		writer    *httptest.ResponseRecorder
		context   stack.Context
		logWriter *bytes.Buffer
	)

	BeforeEach(func() {
		var err error
		request, err = http.NewRequest("GET", "/some/path", nil)
		if err != nil {
			panic(err)
		}

		logWriter = &bytes.Buffer{}

		writer = httptest.NewRecorder()
		ware = middleware.NewRequestLogging(logWriter)
		context = stack.NewContext()
	})

	It("logs the request", func() {
		result := ware.ServeHTTP(writer, request, context)
		Expect(result).To(BeTrue())
		Expect(logWriter.String()).To(Equal("[WEB] request-id: UNKNOWN | GET /some/path\n"))
	})

	It("generates a logger with a prefix that includes the vcap_request_id", func() {
		request.Header.Set("X-Vcap-Request-Id", "some-request-id")

		result := ware.ServeHTTP(writer, request, context)
		Expect(result).To(BeTrue())

		Expect(context.Get("logger")).To(BeAssignableToTypeOf(&log.Logger{}))
		logger := context.Get("logger").(*log.Logger)
		Expect(logger.Prefix()).To(Equal("[WEB] request-id: some-request-id | "))
	})

	Context("when the request id is unknown", func() {
		It("generates a logger with a prefix that states the request id is unknown", func() {
			result := ware.ServeHTTP(writer, request, context)
			Expect(result).To(BeTrue())

			Expect(context.Get("logger")).To(BeAssignableToTypeOf(&log.Logger{}))
			logger := context.Get("logger").(*log.Logger)
			Expect(logger.Prefix()).To(Equal("[WEB] request-id: UNKNOWN | "))
		})
	})
})
