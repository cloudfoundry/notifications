package clients

import (
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/cloudfoundry-incubator/notifications/v1/web/middleware"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

type muxer interface {
	Handle(method, path string, handler stack.Handler, middleware ...stack.Middleware)
	GetRouter() *mux.Router
}

type Routes struct {
	RequestLogging                   stack.Middleware
	NotificationsManageAuthenticator stack.Middleware
	DatabaseAllocator                stack.Middleware

	ErrorWriter      errorWriter
	TemplateAssigner services.TemplateAssignerInterface
}

func (r Routes) Register(m muxer) {
	requestCounter := middleware.NewRequestCounter(m.GetRouter(), metrics.DefaultLogger)
	m.Handle("PUT", "/clients/{client_id}/template", NewAssignTemplateHandler(r.TemplateAssigner, r.ErrorWriter), r.RequestLogging, requestCounter, r.NotificationsManageAuthenticator, r.DatabaseAllocator)
}
