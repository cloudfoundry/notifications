package messages

import "github.com/ryanmoran/stack"

type muxer interface {
	Handle(method, path string, handler stack.Handler, middleware ...stack.Middleware)
}

type Routes struct {
	RequestCounter                               stack.Middleware
	RequestLogging                               stack.Middleware
	NotificationsWriteOrEmailsWriteAuthenticator stack.Middleware
	DatabaseAllocator                            stack.Middleware

	MessageFinder messageFinder
	ErrorWriter   errorWriter
}

func (r Routes) Register(m muxer) {
	m.Handle("GET", "/messages/{message_id}", NewGetHandler(r.MessageFinder, r.ErrorWriter), r.RequestLogging, r.RequestCounter, r.NotificationsWriteOrEmailsWriteAuthenticator, r.DatabaseAllocator)
}
