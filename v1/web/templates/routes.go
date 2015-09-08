package templates

import "github.com/ryanmoran/stack"

type muxer interface {
	Handle(method, path string, handler stack.Handler, middleware ...stack.Middleware)
}

type Routes struct {
	RequestCounter                          stack.Middleware
	RequestLogging                          stack.Middleware
	DatabaseAllocator                       stack.Middleware
	NotificationTemplatesReadAuthenticator  stack.Middleware
	NotificationTemplatesWriteAuthenticator stack.Middleware
	NotificationsManageAuthenticator        stack.Middleware

	ErrorWriter               errorWriter
	TemplateFinder            templateFinder
	TemplateLister            templateLister
	TemplateUpdater           templateUpdater
	TemplateCreator           templateCreator
	TemplateDeleter           templateDeleter
	TemplateAssociationLister templateAssociationLister
}

func (r Routes) Register(m muxer) {
	m.Handle("GET", "/default_template", NewGetDefaultHandler(r.TemplateFinder, r.ErrorWriter), r.RequestLogging, r.RequestCounter, r.NotificationTemplatesReadAuthenticator, r.DatabaseAllocator)
	m.Handle("PUT", "/default_template", NewUpdateDefaultHandler(r.TemplateUpdater, r.ErrorWriter), r.RequestLogging, r.RequestCounter, r.NotificationTemplatesWriteAuthenticator, r.DatabaseAllocator)
	m.Handle("GET", "/templates", NewListHandler(r.TemplateLister, r.ErrorWriter), r.RequestLogging, r.RequestCounter, r.NotificationTemplatesReadAuthenticator, r.DatabaseAllocator)
	m.Handle("POST", "/templates", NewCreateHandler(r.TemplateCreator, r.ErrorWriter), r.RequestLogging, r.RequestCounter, r.NotificationTemplatesWriteAuthenticator, r.DatabaseAllocator)
	m.Handle("GET", "/templates/{template_id}", NewGetHandler(r.TemplateFinder, r.ErrorWriter), r.RequestLogging, r.RequestCounter, r.NotificationTemplatesReadAuthenticator, r.DatabaseAllocator)
	m.Handle("PUT", "/templates/{template_id}", NewUpdateHandler(r.TemplateUpdater, r.ErrorWriter), r.RequestLogging, r.RequestCounter, r.NotificationTemplatesWriteAuthenticator, r.DatabaseAllocator)
	m.Handle("DELETE", "/templates/{template_id}", NewDeleteHandler(r.TemplateDeleter, r.ErrorWriter), r.RequestLogging, r.RequestCounter, r.NotificationTemplatesWriteAuthenticator, r.DatabaseAllocator)
	m.Handle("GET", "/templates/{template_id}/associations", NewListAssociationsHandler(r.TemplateAssociationLister, r.ErrorWriter), r.RequestLogging, r.RequestCounter, r.NotificationsManageAuthenticator, r.DatabaseAllocator)
}
