package senders

import (
	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

type Routes struct {
	RequestLogging    stack.Middleware
	Authenticator     stack.Middleware
	DatabaseAllocator stack.Middleware
	SendersCollection collections.SendersCollection
}

func (r Routes) Register(router *mux.Router) {
	createStack := stack.NewStack(NewCreateHandler(r.SendersCollection)).Use(r.RequestLogging, r.Authenticator, r.DatabaseAllocator)
	getStack := stack.NewStack(NewGetHandler(r.SendersCollection)).Use(r.RequestLogging, r.Authenticator, r.DatabaseAllocator)

	router.Handle("/senders", createStack).Methods("POST").Name("POST /senders")
	router.Handle("/senders/{sender_id}", getStack).Methods("GET").Name("GET /senders/{sender_id}")
}
