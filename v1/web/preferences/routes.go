package preferences

import (
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/ryanmoran/stack"
)

type muxer interface {
	Handle(method, path string, handler stack.Handler, middleware ...stack.Middleware)
}

type preferenceUpdater interface {
	Update(connection services.ConnectionInterface, preferences []models.Preference, globallyUnsubscribe bool, userID string) error
}

type Routes struct {
	CORS                                      stack.Middleware
	RequestCounter                            stack.Middleware
	RequestLogging                            stack.Middleware
	DatabaseAllocator                         stack.Middleware
	NotificationPreferencesReadAuthenticator  stack.Middleware
	NotificationPreferencesAdminAuthenticator stack.Middleware
	NotificationPreferencesWriteAuthenticator stack.Middleware

	ErrorWriter       errorWriter
	PreferencesFinder preferencesFinder
	PreferenceUpdater preferenceUpdater
}

func (r Routes) Register(m muxer) {
	m.Handle("OPTIONS", "/user_preferences", NewOptionsHandler(), r.RequestLogging, r.RequestCounter, r.CORS)
	m.Handle("OPTIONS", "/user_preferences/{user_id}", NewOptionsHandler(), r.RequestLogging, r.RequestCounter, r.CORS)
	m.Handle("GET", "/user_preferences", NewGetPreferencesHandler(r.PreferencesFinder, r.ErrorWriter), r.RequestLogging, r.RequestCounter, r.CORS, r.NotificationPreferencesReadAuthenticator, r.DatabaseAllocator)
	m.Handle("PATCH", "/user_preferences", NewUpdatePreferencesHandler(r.PreferenceUpdater, r.ErrorWriter), r.RequestLogging, r.RequestCounter, r.CORS, r.NotificationPreferencesWriteAuthenticator, r.DatabaseAllocator)
	m.Handle("GET", "/user_preferences/{user_id}", NewGetUserPreferencesHandler(r.PreferencesFinder, r.ErrorWriter), r.RequestLogging, r.RequestCounter, r.CORS, r.NotificationPreferencesAdminAuthenticator, r.DatabaseAllocator)
	m.Handle("PATCH", "/user_preferences/{user_id}", NewUpdateUserPreferencesHandler(r.PreferenceUpdater, r.ErrorWriter), r.RequestLogging, r.RequestCounter, r.CORS, r.NotificationPreferencesAdminAuthenticator, r.DatabaseAllocator)
}
