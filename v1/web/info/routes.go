package info

import (
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/v1/web/middleware"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

type muxer interface {
	Handle(method, path string, handler stack.Handler, middleware ...stack.Middleware)
	GetRouter() *mux.Router
}

type Routes struct {
	RequestLogging stack.Middleware
}

func (r Routes) Register(m muxer) {
	requestCounter := middleware.NewRequestCounter(m.GetRouter(), metrics.DefaultLogger)
	m.Handle("GET", "/info", NewGetHandler(), r.RequestLogging, requestCounter)
}
