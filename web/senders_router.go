package web

import (
	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/web/senders"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

func NewSendersRouter(requestLogging RequestLogging, authenticator Authenticator, databaseAllocator DatabaseAllocator, sendersCollection collections.SendersCollection) *mux.Router {
	router := mux.NewRouter()

	createStack := stack.NewStack(senders.NewCreateHandler(sendersCollection)).Use(requestLogging, authenticator, databaseAllocator)
	getStack := stack.NewStack(senders.NewGetHandler(sendersCollection)).Use(requestLogging, authenticator, databaseAllocator)

	router.Handle("/senders", createStack).Methods("POST").Name("POST /senders")
	router.Handle("/senders/{sender_id}", getStack).Methods("GET").Name("GET /senders/{sender_id}")

	return router
}
