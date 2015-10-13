package unsubscribers

import (
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/ryanmoran/stack"
)

type muxer interface {
	Handle(method, path string, handler stack.Handler, middleware ...stack.Middleware)
}

type Routes struct {
	RequestLogging          stack.Middleware
	Authenticator           stack.Middleware
	DatabaseAllocator       stack.Middleware
	UnsubscribersCollection collections.UnsubscribersCollection
}

func (r Routes) Register(m muxer) {
	m.Handle("PUT", "/campaign_types/{campaign_type_id}/unsubscribers/{user_guid}", NewUpdateHandler(r.UnsubscribersCollection), r.RequestLogging, r.Authenticator, r.DatabaseAllocator)
	m.Handle("DELETE", "/campaign_types/{campaign_type_id}/unsubscribers/{user_guid}", NewDeleteHandler(r.UnsubscribersCollection), r.RequestLogging, r.Authenticator, r.DatabaseAllocator)
}
