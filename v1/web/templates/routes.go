package templates

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
	RequestLogging                          stack.Middleware
	DatabaseAllocator                       stack.Middleware
	NotificationTemplatesReadAuthenticator  stack.Middleware
	NotificationTemplatesWriteAuthenticator stack.Middleware
	NotificationsManageAuthenticator        stack.Middleware

	ErrorWriter               errorWriter
	TemplateFinder            services.TemplateFinderInterface
	TemplateLister            services.TemplateListerInterface
	TemplateUpdater           services.TemplateUpdaterInterface
	TemplateCreator           services.TemplateCreatorInterface
	TemplateDeleter           services.TemplateDeleterInterface
	TemplateAssociationLister services.TemplateAssociationListerInterface
}

func (r Routes) Register(m muxer) {
	requestCounter := middleware.NewRequestCounter(m.GetRouter(), metrics.DefaultLogger)
	m.Handle("GET", "/default_template", NewGetDefaultHandler(r.TemplateFinder, r.ErrorWriter), r.RequestLogging, requestCounter, r.NotificationTemplatesReadAuthenticator, r.DatabaseAllocator)
	m.Handle("PUT", "/default_template", NewUpdateDefaultHandler(r.TemplateUpdater, r.ErrorWriter), r.RequestLogging, requestCounter, r.NotificationTemplatesWriteAuthenticator, r.DatabaseAllocator)
	m.Handle("GET", "/templates", NewListHandler(r.TemplateLister, r.ErrorWriter), r.RequestLogging, requestCounter, r.NotificationTemplatesReadAuthenticator, r.DatabaseAllocator)
	m.Handle("POST", "/templates", NewCreateHandler(r.TemplateCreator, r.ErrorWriter), r.RequestLogging, requestCounter, r.NotificationTemplatesWriteAuthenticator, r.DatabaseAllocator)
	m.Handle("GET", "/templates/{template_id}", NewGetHandler(r.TemplateFinder, r.ErrorWriter), r.RequestLogging, requestCounter, r.NotificationTemplatesReadAuthenticator, r.DatabaseAllocator)
	m.Handle("PUT", "/templates/{template_id}", NewUpdateHandler(r.TemplateUpdater, r.ErrorWriter), r.RequestLogging, requestCounter, r.NotificationTemplatesWriteAuthenticator, r.DatabaseAllocator)
	m.Handle("DELETE", "/templates/{template_id}", NewDeleteHandler(r.TemplateDeleter, r.ErrorWriter), r.RequestLogging, requestCounter, r.NotificationTemplatesWriteAuthenticator, r.DatabaseAllocator)
	m.Handle("GET", "/templates/{template_id}/associations", NewListAssociationsHandler(r.TemplateAssociationLister, r.ErrorWriter), r.RequestLogging, requestCounter, r.NotificationsManageAuthenticator, r.DatabaseAllocator)
}
