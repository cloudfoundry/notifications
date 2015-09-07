package clients

import "github.com/ryanmoran/stack"

type muxer interface {
	Handle(method, path string, handler stack.Handler, middleware ...stack.Middleware)
}

type Routes struct {
	RequestCounter                   stack.Middleware
	RequestLogging                   stack.Middleware
	NotificationsManageAuthenticator stack.Middleware
	DatabaseAllocator                stack.Middleware

	ErrorWriter      errorWriter
	TemplateAssigner assignsTemplates
}

func (r Routes) Register(m muxer) {
	m.Handle("PUT", "/clients/{client_id}/template", NewAssignTemplateHandler(r.TemplateAssigner, r.ErrorWriter), r.RequestLogging, r.RequestCounter, r.NotificationsManageAuthenticator, r.DatabaseAllocator)
}
