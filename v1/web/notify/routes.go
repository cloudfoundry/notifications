package notify

import "github.com/ryanmoran/stack"

type muxer interface {
	Handle(method, path string, handler stack.Handler, middleware ...stack.Middleware)
}

type Routes struct {
	RequestCounter                  stack.Middleware
	RequestLogging                  stack.Middleware
	DatabaseAllocator               stack.Middleware
	NotificationsWriteAuthenticator stack.Middleware
	EmailsWriteAuthenticator        stack.Middleware

	Notify               notifyExecutor
	ErrorWriter          errorWriter
	UserStrategy         Dispatcher
	SpaceStrategy        Dispatcher
	OrganizationStrategy Dispatcher
	EveryoneStrategy     Dispatcher
	UAAScopeStrategy     Dispatcher
	EmailStrategy        Dispatcher
}

func (r Routes) Register(m muxer) {
	m.Handle("POST", "/users/{user_id}", NewUserHandler(r.Notify, r.ErrorWriter, r.UserStrategy), r.RequestLogging, r.RequestCounter, r.NotificationsWriteAuthenticator, r.DatabaseAllocator)
	m.Handle("POST", "/spaces/{space_id}", NewSpaceHandler(r.Notify, r.ErrorWriter, r.SpaceStrategy), r.RequestLogging, r.RequestCounter, r.NotificationsWriteAuthenticator, r.DatabaseAllocator)
	m.Handle("POST", "/organizations/{org_id}", NewOrganizationHandler(r.Notify, r.ErrorWriter, r.OrganizationStrategy), r.RequestLogging, r.RequestCounter, r.NotificationsWriteAuthenticator, r.DatabaseAllocator)
	m.Handle("POST", "/everyone", NewEveryoneHandler(r.Notify, r.ErrorWriter, r.EveryoneStrategy), r.RequestLogging, r.RequestCounter, r.NotificationsWriteAuthenticator, r.DatabaseAllocator)
	m.Handle("POST", "/uaa_scopes/{scope}", NewUAAScopeHandler(r.Notify, r.ErrorWriter, r.UAAScopeStrategy), r.RequestLogging, r.RequestCounter, r.NotificationsWriteAuthenticator, r.DatabaseAllocator)
	m.Handle("POST", "/emails", NewEmailHandler(r.Notify, r.ErrorWriter, r.EmailStrategy), r.RequestLogging, r.RequestCounter, r.EmailsWriteAuthenticator, r.DatabaseAllocator)
}
