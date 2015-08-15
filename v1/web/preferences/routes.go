package preferences

import (
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/cloudfoundry-incubator/notifications/v1/web/middleware"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

type muxer interface {
	Handle(method, path string, handler stack.Handler, middleware ...stack.Middleware)
	GetRouter() *mux.Router
}

type Routes struct {
	CORS                                      stack.Middleware
	RequestLogging                            stack.Middleware
	DatabaseAllocator                         stack.Middleware
	NotificationPreferencesReadAuthenticator  stack.Middleware
	NotificationPreferencesAdminAuthenticator stack.Middleware
	NotificationPreferencesWriteAuthenticator stack.Middleware

	ErrorWriter       errorWriter
	PreferencesFinder services.PreferencesFinderInterface
	PreferenceUpdater services.PreferenceUpdaterInterface
}

func (r Routes) Register(m muxer) {
	requestCounter := middleware.NewRequestCounter(m.GetRouter(), metrics.DefaultLogger)

	m.Handle("OPTIONS", "/user_preferences", NewOptionsHandler(), r.RequestLogging, requestCounter, r.CORS)
	m.Handle("OPTIONS", "/user_preferences/{user_id}", NewOptionsHandler(), r.RequestLogging, requestCounter, r.CORS)
	m.Handle("GET", "/user_preferences", NewGetPreferencesHandler(r.PreferencesFinder, r.ErrorWriter), r.RequestLogging, requestCounter, r.CORS, r.NotificationPreferencesReadAuthenticator, r.DatabaseAllocator)
	m.Handle("PATCH", "/user_preferences", NewUpdatePreferencesHandler(r.PreferenceUpdater, r.ErrorWriter), r.RequestLogging, requestCounter, r.CORS, r.NotificationPreferencesWriteAuthenticator, r.DatabaseAllocator)
	m.Handle("GET", "/user_preferences/{user_id}", NewGetUserPreferencesHandler(r.PreferencesFinder, r.ErrorWriter), r.RequestLogging, requestCounter, r.CORS, r.NotificationPreferencesAdminAuthenticator, r.DatabaseAllocator)
	m.Handle("PATCH", "/user_preferences/{user_id}", NewUpdateUserPreferencesHandler(r.PreferenceUpdater, r.ErrorWriter), r.RequestLogging, requestCounter, r.CORS, r.NotificationPreferencesAdminAuthenticator, r.DatabaseAllocator)
}
