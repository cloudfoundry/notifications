package web

import (
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

func NewMessagesRouter(messageFinder services.MessageFinderInterface,
	errorWriter handlers.ErrorWriterInterface,
	logging RequestLogging,
	notificationsWriteOrEmailsWriteAuthenticator Authenticator,
	databaseAllocator DatabaseAllocator) *mux.Router {

	router := mux.NewRouter()
	requestCounter := NewRequestCounter(router, metrics.DefaultLogger)

	router.Handle("/messages/{message_id}", stack.NewStack(handlers.NewGetMessages(messageFinder, errorWriter)).Use(logging, requestCounter, notificationsWriteOrEmailsWriteAuthenticator, databaseAllocator)).Methods("GET").Name("GET /messages/{message_id}")

	return router
}
