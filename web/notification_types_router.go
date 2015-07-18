package web

import (
	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/web/notificationtypes"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

type NotificationTypesRouterConfig struct {
	NotificationTypesCollection collections.NotificationTypesCollection
}

func NewNotificationTypesRouter(config NotificationTypesRouterConfig) *mux.Router {
	router := mux.NewRouter()

	createStack := stack.NewStack(notificationtypes.NewCreateHandler(nil))

	router.Handle("/senders/{sender_id}/notification_types", createStack).Methods("POST").Name("POST /senders/{sender_id}/notification_types")

	return router
}
