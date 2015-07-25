package info

import (
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

type Routes struct {
	Version        int
	RequestLogging stack.Middleware
}

func (r Routes) Register(router *mux.Router) {
	requestCounter := middleware.NewRequestCounter(router, metrics.DefaultLogger)

	router.Handle("/info", stack.NewStack(NewGetHandler(r.Version)).Use(r.RequestLogging, requestCounter)).Methods("GET").Name("GET /info")
}
