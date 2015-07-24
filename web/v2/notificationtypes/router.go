package notificationtypes

import (
	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

type RouterConfig struct {
	RequestLogging              stack.Middleware
	Authenticator               middleware.Authenticator
	DatabaseAllocator           middleware.DatabaseAllocator
	NotificationTypesCollection collections.NotificationTypesCollection
}

func NewRouter(config RouterConfig) *mux.Router {
	router := mux.NewRouter()

	createStack := stack.NewStack(NewCreateHandler(config.NotificationTypesCollection)).Use(config.RequestLogging, config.Authenticator, config.DatabaseAllocator)
	listStack := stack.NewStack(NewListHandler(config.NotificationTypesCollection)).Use(config.RequestLogging, config.Authenticator, config.DatabaseAllocator)
	showStack := stack.NewStack(NewShowHandler(config.NotificationTypesCollection)).Use(config.RequestLogging, config.Authenticator, config.DatabaseAllocator)

	router.Handle("/senders/{sender_id}/notification_types", createStack).Methods("POST").Name("POST /senders/{sender_id}/notification_types")
	router.Handle("/senders/{sender_id}/notification_types/{notification_type_id}", showStack).Methods("GET").Name("GET /senders/{sender_id}/notification_types/{notification_type_id}")
	router.Handle("/senders/{sender_id}/notification_types", listStack).Methods("GET").Name("GET /senders/{sender_id}/notification_types")

	return router
}
