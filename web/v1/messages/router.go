package messages

import (
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

type RouterConfig struct {
	RequestLogging                               stack.Middleware
	NotificationsWriteOrEmailsWriteAuthenticator stack.Middleware
	DatabaseAllocator                            stack.Middleware

	MessageFinder services.MessageFinderInterface
	ErrorWriter   errorWriter
}

func NewRouter(config RouterConfig) *mux.Router {
	router := mux.NewRouter()
	requestCounter := middleware.NewRequestCounter(router, metrics.DefaultLogger)

	getHandler := NewGetHandler(config.MessageFinder, config.ErrorWriter)
	getStack := stack.NewStack(getHandler).Use(config.RequestLogging, requestCounter, config.NotificationsWriteOrEmailsWriteAuthenticator, config.DatabaseAllocator)
	router.Handle("/messages/{message_id}", getStack).Methods("GET").Name("GET /messages/{message_id}")

	return router
}
