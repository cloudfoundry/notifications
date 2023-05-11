package middleware_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/web/middleware"
	"github.com/pivotal-golang/lager"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type logLine struct {
	Source   string
	Message  string
	LogLevel int
	Data     map[string]interface{}
}

var _ = Describe("RequestLogging", func() {
	var (
		ware      middleware.RequestLogging
		request   *http.Request
		writer    *httptest.ResponseRecorder
		context   stack.Context
		logger    lager.Logger
		logWriter *bytes.Buffer
		clock     *mocks.Clock
	)

	BeforeEach(func() {
		var err error
		request, err = http.NewRequest("GET", "/some/path", nil)
		Expect(err).NotTo(HaveOccurred())

		request.Header.Set("X-Vcap-Request-Id", "some-request-id")

		logWriter = &bytes.Buffer{}
		logger = lager.NewLogger("my-app")
		logger.RegisterSink(lager.NewWriterSink(logWriter, lager.DEBUG))

		writer = httptest.NewRecorder()
		context = stack.NewContext()
		clock = mocks.NewClock()
		clock.NowCall.Returns.Time = time.Now()

		ware = middleware.NewRequestLogging(logger, clock)
	})

	It("logs the request without an API version", func() {
		result := ware.ServeHTTP(writer, request, context)
		Expect(result).To(BeTrue())

		var line logLine
		err := json.Unmarshal(logWriter.Bytes(), &line)
		Expect(err).NotTo(HaveOccurred())
		Expect(line.Source).To(Equal("my-app"))
		Expect(line.Message).To(Equal("my-app.request.incoming"))
		Expect(line.LogLevel).To(Equal(int(lager.DEBUG)))
		Expect(line.Data).To(HaveKeyWithValue("vcap_request_id", "some-request-id"))
		Expect(line.Data).NotTo(HaveKey("api_version"))
		Expect(line.Data).To(HaveKeyWithValue("method", "GET"))
		Expect(line.Data).To(HaveKeyWithValue("path", "/some/path"))
	})

	It("logs the request with an API version", func() {
		request.Header.Set("X-NOTIFICATIONS-VERSION", "1")

		result := ware.ServeHTTP(writer, request, context)
		Expect(result).To(BeTrue())

		var line logLine
		err := json.Unmarshal(logWriter.Bytes(), &line)
		Expect(err).NotTo(HaveOccurred())
		Expect(line.Source).To(Equal("my-app"))
		Expect(line.Message).To(Equal("my-app.request.incoming"))
		Expect(line.LogLevel).To(Equal(int(lager.DEBUG)))
		Expect(line.Data).To(HaveKeyWithValue("vcap_request_id", "some-request-id"))
		Expect(line.Data).To(HaveKeyWithValue("api_version", "1"))
		Expect(line.Data).To(HaveKeyWithValue("method", "GET"))
		Expect(line.Data).To(HaveKeyWithValue("path", "/some/path"))
	})

	It("adds a logger to the context that includes the vcap_request_id", func() {
		result := ware.ServeHTTP(writer, request, context)
		Expect(result).To(BeTrue())

		logger := context.Get("logger").(lager.Logger)
		logger.Info("hello")

		lines := bytes.Split(logWriter.Bytes(), []byte("\n"))

		var line logLine
		err := json.Unmarshal(lines[1], &line)
		Expect(err).NotTo(HaveOccurred())
		Expect(line.Source).To(Equal("my-app"))
		Expect(line.Message).To(Equal("my-app.request.hello"))
		Expect(line.LogLevel).To(Equal(int(lager.DEBUG)))
		Expect(line.Data).To(HaveKeyWithValue("vcap_request_id", "some-request-id"))
	})

	It("adds the request id to the context", func() {
		result := ware.ServeHTTP(writer, request, context)
		Expect(result).To(BeTrue())

		requestID, ok := context.Get(middleware.VCAPRequestIDKey).(string)
		Expect(ok).To(BeTrue())
		Expect(requestID).To(Equal("some-request-id"))
	})

	It("adds the current time to the context", func() {
		now := time.Now()
		clock.NowCall.Returns.Time = now

		result := ware.ServeHTTP(writer, request, context)
		Expect(result).To(BeTrue())

		requestReceivedTime, ok := context.Get(middleware.RequestReceivedTime).(time.Time)
		Expect(ok).To(BeTrue())
		Expect(requestReceivedTime).To(Equal(now.UTC()))
	})

	Context("when the request id is unknown", func() {
		It("generates a logger with a prefix that states the request id is unknown", func() {
			request.Header.Del("X-Vcap-Request-Id")

			result := ware.ServeHTTP(writer, request, context)
			Expect(result).To(BeTrue())

			logger := context.Get("logger").(lager.Logger)
			logger.Info("hello")

			lines := bytes.Split(logWriter.Bytes(), []byte("\n"))

			var line logLine
			err := json.Unmarshal(lines[1], &line)
			Expect(err).NotTo(HaveOccurred())
			Expect(line.Source).To(Equal("my-app"))
			Expect(line.Message).To(Equal("my-app.request.hello"))
			Expect(line.LogLevel).To(Equal(int(lager.DEBUG)))
			Expect(line.Data).To(HaveKeyWithValue("vcap_request_id", "UNKNOWN"))
		})
	})
})
