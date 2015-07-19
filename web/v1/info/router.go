package info

import (
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

type RouterConfig struct {
	Version        int
	RequestLogging stack.Middleware
}

func NewRouter(config RouterConfig) *mux.Router {
	router := mux.NewRouter()
	requestCounter := middleware.NewRequestCounter(router, metrics.DefaultLogger)

	router.Handle("/info", stack.NewStack(handlers.NewGetInfo(config.Version)).Use(config.RequestLogging, requestCounter)).Methods("GET").Name("GET /info")

	return router
}
