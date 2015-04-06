package middleware

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/ryanmoran/stack"
)

type RequestLogging struct {
	logWriter io.Writer
}

func NewRequestLogging(logWriter io.Writer) RequestLogging {
	return RequestLogging{
		logWriter: logWriter,
	}
}

func (r RequestLogging) ServeHTTP(response http.ResponseWriter, request *http.Request, context stack.Context) bool {
	requestID := request.Header.Get("X-Vcap-Request-Id")
	if requestID == "" {
		requestID = "UNKNOWN"
	}

	logPrefix := fmt.Sprintf("[WEB] request-id: %s | ", requestID)
	logger := log.New(r.logWriter, logPrefix, 0)
	logger.Printf("%s %s\n", request.Method, request.URL.Path)

	context.Set("logger", logger)

	return true
}
