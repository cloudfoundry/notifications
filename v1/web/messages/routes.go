package messages

import (
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

type muxer interface {
	Handle(method, path string, handler stack.Handler, middleware ...stack.Middleware)
	GetRouter() *mux.Router
}

type Routes struct {
	RequestLogging                               stack.Middleware
	NotificationsWriteOrEmailsWriteAuthenticator stack.Middleware
	DatabaseAllocator                            stack.Middleware

	MessageFinder services.MessageFinderInterface
	ErrorWriter   errorWriter
}

func (r Routes) Register(m muxer) {
	requestCounter := middleware.NewRequestCounter(m.GetRouter(), metrics.DefaultLogger)
	m.Handle("GET", "/messages/{message_id}", NewGetHandler(r.MessageFinder, r.ErrorWriter), r.RequestLogging, requestCounter, r.NotificationsWriteOrEmailsWriteAuthenticator, r.DatabaseAllocator)
}
