package notifications

import "github.com/ryanmoran/stack"

type muxer interface {
	Handle(method, path string, handler stack.Handler, middleware ...stack.Middleware)
}

type Routes struct {
	RequestCounter                   stack.Middleware
	RequestLogging                   stack.Middleware
	DatabaseAllocator                stack.Middleware
	NotificationsWriteAuthenticator  stack.Middleware
	NotificationsManageAuthenticator stack.Middleware

	ErrorWriter          errorWriter
	Registrar            registrar
	TemplateAssigner     assignsTemplates
	NotificationsFinder  listsAllClientsAndNotifications
	NotificationsUpdater notificationsUpdater
}

func (r Routes) Register(m muxer) {
	m.Handle("PUT", "/registration", NewRegistrationHandler(r.Registrar, r.ErrorWriter), r.RequestLogging, r.RequestCounter, r.NotificationsWriteAuthenticator, r.DatabaseAllocator)
	m.Handle("PUT", "/notifications", NewPutHandler(r.Registrar, r.ErrorWriter), r.RequestLogging, r.RequestCounter, r.NotificationsWriteAuthenticator, r.DatabaseAllocator)
	m.Handle("GET", "/notifications", NewListHandler(r.NotificationsFinder, r.ErrorWriter), r.RequestLogging, r.RequestCounter, r.NotificationsManageAuthenticator, r.DatabaseAllocator)
	m.Handle("PUT", "/clients/{client_id}/notifications/{notification_id}", NewUpdateHandler(r.NotificationsUpdater, r.ErrorWriter), r.RequestLogging, r.RequestCounter, r.NotificationsManageAuthenticator, r.DatabaseAllocator)
	m.Handle("PUT", "/clients/{client_id}/notifications/{notification_id}/template", NewAssignTemplateHandler(r.TemplateAssigner, r.ErrorWriter), r.RequestLogging, r.RequestCounter, r.NotificationsManageAuthenticator, r.DatabaseAllocator)
}
