package notificationtypes

import (
	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

type Routes struct {
	RequestLogging              stack.Middleware
	Authenticator               middleware.Authenticator
	DatabaseAllocator           middleware.DatabaseAllocator
	NotificationTypesCollection collections.NotificationTypesCollection
}

func (r Routes) Register(router *mux.Router) {
	createStack := stack.NewStack(NewCreateHandler(r.NotificationTypesCollection)).Use(r.RequestLogging, r.Authenticator, r.DatabaseAllocator)
	listStack := stack.NewStack(NewListHandler(r.NotificationTypesCollection)).Use(r.RequestLogging, r.Authenticator, r.DatabaseAllocator)

	router.Handle("/senders/{sender_id}/notification_types", createStack).Methods("POST").Name("POST /senders/{sender_id}/notification_types")
	router.Handle("/senders/{sender_id}/notification_types", listStack).Methods("GET").Name("GET /senders/{sender_id}/notification_types")
}
