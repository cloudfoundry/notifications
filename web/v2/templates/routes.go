package templates

import (
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/ryanmoran/stack"
)

type muxer interface {
	Handle(method, path string, handler stack.Handler, middleware ...stack.Middleware)
}

type Routes struct {
	RequestLogging      stack.Middleware
	Authenticator       middleware.Authenticator
	DatabaseAllocator   middleware.DatabaseAllocator
	TemplatesCollection interface{}
}

func (r Routes) Register(m muxer) {
	m.Handle("POST", "/templates", NewCreateHandler(r.TemplatesCollection), r.RequestLogging, r.Authenticator, r.DatabaseAllocator)
}
