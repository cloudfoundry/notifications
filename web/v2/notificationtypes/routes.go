package notificationtypes

import (
	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/ryanmoran/stack"
)

type muxer interface {
	Handle(method, path string, handler stack.Handler, middleware ...stack.Middleware)
}

type Routes struct {
	RequestLogging              stack.Middleware
	Authenticator               middleware.Authenticator
	DatabaseAllocator           middleware.DatabaseAllocator
	NotificationTypesCollection collections.NotificationTypesCollection
}

func (r Routes) Register(m muxer) {
	m.Handle("POST", "/senders/{sender_id}/notification_types", NewCreateHandler(r.NotificationTypesCollection), r.RequestLogging, r.Authenticator, r.DatabaseAllocator)
	m.Handle("GET", "/senders/{sender_id}/notification_types", NewListHandler(r.NotificationTypesCollection), r.RequestLogging, r.Authenticator, r.DatabaseAllocator)
}
