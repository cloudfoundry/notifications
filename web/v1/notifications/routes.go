package notifications

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
	RequestLogging                   stack.Middleware
	DatabaseAllocator                stack.Middleware
	NotificationsWriteAuthenticator  stack.Middleware
	NotificationsManageAuthenticator stack.Middleware

	ErrorWriter          errorWriter
	Registrar            services.RegistrarInterface
	TemplateAssigner     services.TemplateAssignerInterface
	NotificationsFinder  services.NotificationsFinderInterface
	NotificationsUpdater services.NotificationsUpdaterInterface
}

func (r Routes) Register(m muxer) {
	requestCounter := middleware.NewRequestCounter(m.GetRouter(), metrics.DefaultLogger)
	m.Handle("PUT", "/registration", NewRegistrationHandler(r.Registrar, r.ErrorWriter), r.RequestLogging, requestCounter, r.NotificationsWriteAuthenticator, r.DatabaseAllocator)
	m.Handle("PUT", "/notifications", NewPutHandler(r.Registrar, r.ErrorWriter), r.RequestLogging, requestCounter, r.NotificationsWriteAuthenticator, r.DatabaseAllocator)
	m.Handle("GET", "/notifications", NewListHandler(r.NotificationsFinder, r.ErrorWriter), r.RequestLogging, requestCounter, r.NotificationsManageAuthenticator, r.DatabaseAllocator)
	m.Handle("PUT", "/clients/{client_id}/notifications/{notification_id}", NewUpdateHandler(r.NotificationsUpdater, r.ErrorWriter), r.RequestLogging, requestCounter, r.NotificationsManageAuthenticator, r.DatabaseAllocator)
	m.Handle("PUT", "/clients/{client_id}/notifications/{notification_id}/template", NewAssignTemplateHandler(r.TemplateAssigner, r.ErrorWriter), r.RequestLogging, requestCounter, r.NotificationsManageAuthenticator, r.DatabaseAllocator)
}
