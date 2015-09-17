package middleware

import (
	"net/http"
	"time"

	"github.com/pivotal-golang/lager"
	"github.com/ryanmoran/stack"
)

const (
	VCAPRequestIDKey    = "vcap_request_id"
	APIVersion          = "api_version"
	RequestReceivedTime = "request_received_time"
)

type RequestLogging struct {
	logger lager.Logger
}

func NewRequestLogging(logger lager.Logger) RequestLogging {
	return RequestLogging{logger}
}

func (r RequestLogging) ServeHTTP(response http.ResponseWriter, request *http.Request, context stack.Context) bool {
	requestID := request.Header.Get("X-Vcap-Request-Id")
	if requestID == "" {
		requestID = "UNKNOWN"
	}

	apiVersion := request.Header.Get("X-NOTIFICATIONS-VERSION")

	logSession := r.logger.Session("request", lager.Data{
		VCAPRequestIDKey: requestID,
		APIVersion:       apiVersion,
	})

	logSession.Info("incoming", lager.Data{
		"method": request.Method,
		"path":   request.URL.Path,
	})

	context.Set("logger", logSession)
	context.Set(VCAPRequestIDKey, requestID)
	context.Set(RequestReceivedTime, time.Now().UTC())

	return true
}
