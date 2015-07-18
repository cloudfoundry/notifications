package web

import (
	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/web/senders"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

type SendersRouterConfig struct {
	RequestLogging    RequestLogging
	Authenticator     Authenticator
	DatabaseAllocator DatabaseAllocator
	SendersCollection collections.SendersCollection
}

func NewSendersRouter(config SendersRouterConfig) *mux.Router {
	router := mux.NewRouter()

	createStack := stack.NewStack(senders.NewCreateHandler(config.SendersCollection)).Use(config.RequestLogging, config.Authenticator, config.DatabaseAllocator)
	getStack := stack.NewStack(senders.NewGetHandler(config.SendersCollection)).Use(config.RequestLogging, config.Authenticator, config.DatabaseAllocator)

	router.Handle("/senders", createStack).Methods("POST").Name("POST /senders")
	router.Handle("/senders/{sender_id}", getStack).Methods("GET").Name("GET /senders/{sender_id}")

	return router
}
