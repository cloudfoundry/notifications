package templates

import (
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/ryanmoran/stack"
)

type muxer interface {
	Handle(method, path string, handler stack.Handler, middleware ...stack.Middleware)
}

type Routes struct {
	RequestLogging      stack.Middleware
	WriteAuthenticator  stack.Middleware
	AdminAuthenticator  stack.Middleware
	DatabaseAllocator   stack.Middleware
	TemplatesCollection collections.TemplatesCollection
}

func (r Routes) Register(m muxer) {
	m.Handle("GET", "/templates", NewListHandler(r.TemplatesCollection), r.RequestLogging, r.WriteAuthenticator, r.DatabaseAllocator)
	m.Handle("POST", "/templates", NewCreateHandler(r.TemplatesCollection), r.RequestLogging, r.WriteAuthenticator, r.DatabaseAllocator)
	m.Handle("GET", "/templates/{template_id}", NewGetHandler(r.TemplatesCollection), r.RequestLogging, r.WriteAuthenticator, r.DatabaseAllocator)
	m.Handle("DELETE", "/templates/{template_id}", NewDeleteHandler(r.TemplatesCollection), r.RequestLogging, r.WriteAuthenticator, r.DatabaseAllocator)
	m.Handle("PUT", "/templates/default", NewUpdateDefaultHandler(r.TemplatesCollection), r.RequestLogging, r.AdminAuthenticator, r.DatabaseAllocator)
	m.Handle("PUT", "/templates/{template_id}", NewUpdateHandler(r.TemplatesCollection), r.RequestLogging, r.WriteAuthenticator, r.DatabaseAllocator)
}
