package middleware

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/pivotal-golang/lager"
	"github.com/ryanmoran/stack"
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

	logSession := r.logger.Session("request", lager.Data{
		handlers.VCAPRequestIDKey: requestID,
	})

	logSession.Info("incoming", lager.Data{
		"method": request.Method,
		"path":   request.URL.Path,
	})

	context.Set("logger", logSession)
	context.Set(handlers.VCAPRequestIDKey, requestID)

	return true
}
