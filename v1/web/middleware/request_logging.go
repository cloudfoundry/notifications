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

type clock interface {
	Now() time.Time
}

type RequestLogging struct {
	logger lager.Logger
	clock  clock
}

func NewRequestLogging(logger lager.Logger, clock clock) RequestLogging {
	return RequestLogging{
		logger: logger,
		clock:  clock,
	}
}

func (r RequestLogging) ServeHTTP(response http.ResponseWriter, request *http.Request, context stack.Context) bool {
	requestID := request.Header.Get("X-Vcap-Request-Id")
	if requestID == "" {
		requestID = "UNKNOWN"
	}

	logData := lager.Data{
		VCAPRequestIDKey: requestID,
	}

	apiVersion := request.Header.Get("X-NOTIFICATIONS-VERSION")
	if apiVersion != "" {
		logData[APIVersion] = apiVersion
	}

	logSession := r.logger.Session("request", logData)

	logSession.Info("incoming", lager.Data{
		"method": request.Method,
		"path":   request.URL.Path,
	})

	context.Set("logger", logSession)
	context.Set(VCAPRequestIDKey, requestID)
	context.Set(RequestReceivedTime, r.clock.Now().UTC())

	return true
}
