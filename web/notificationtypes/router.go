package notificationtypes

import (
	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

type RouterConfig struct {
	NotificationTypesCollection collections.NotificationTypesCollection
}

func NewRouter(config RouterConfig) *mux.Router {
	router := mux.NewRouter()

	createStack := stack.NewStack(NewCreateHandler(nil))

	router.Handle("/senders/{sender_id}/notification_types", createStack).Methods("POST").Name("POST /senders/{sender_id}/notification_types")

	return router
}
