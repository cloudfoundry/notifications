package templates

import (
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/web/middleware"
	"github.com/ryanmoran/stack"
)

type muxer interface {
	Handle(method, path string, handler stack.Handler, middleware ...stack.Middleware)
}

type Routes struct {
	RequestLogging      stack.Middleware
	Authenticator       middleware.Authenticator
	DatabaseAllocator   middleware.DatabaseAllocator
	TemplatesCollection collections.TemplatesCollection
}

func (r Routes) Register(m muxer) {
	m.Handle("POST", "/templates", NewCreateHandler(r.TemplatesCollection), r.RequestLogging, r.Authenticator, r.DatabaseAllocator)
	m.Handle("GET", "/templates/{template_id}", NewGetHandler(r.TemplatesCollection), r.RequestLogging, r.Authenticator, r.DatabaseAllocator)
	m.Handle("DELETE", "/templates/{template_id}", NewDeleteHandler(r.TemplatesCollection), r.RequestLogging, r.Authenticator, r.DatabaseAllocator)
}
