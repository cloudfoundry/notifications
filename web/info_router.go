package web

import (
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

func NewInfoRouter(version int, logging RequestLogging) *mux.Router {
	router := mux.NewRouter()
	requestCounter := NewRequestCounter(router, metrics.DefaultLogger)

	router.Handle("/info", stack.NewStack(handlers.NewGetInfo(version)).Use(logging, requestCounter)).Methods("GET").Name("GET /info")

	return router
}
