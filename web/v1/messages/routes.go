package messages

import (
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

type Routes struct {
	RequestLogging                               stack.Middleware
	NotificationsWriteOrEmailsWriteAuthenticator stack.Middleware
	DatabaseAllocator                            stack.Middleware

	MessageFinder services.MessageFinderInterface
	ErrorWriter   errorWriter
}

func (r Routes) Register(router *mux.Router) {
	requestCounter := middleware.NewRequestCounter(router, metrics.DefaultLogger)

	getHandler := NewGetHandler(r.MessageFinder, r.ErrorWriter)
	getStack := stack.NewStack(getHandler).Use(r.RequestLogging, requestCounter, r.NotificationsWriteOrEmailsWriteAuthenticator, r.DatabaseAllocator)
	router.Handle("/messages/{message_id}", getStack).Methods("GET").Name("GET /messages/{message_id}")
}
