package preferences

import (
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

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

func (r Routes) Register(router *mux.Router) {
	requestCounter := middleware.NewRequestCounter(router, metrics.DefaultLogger)

	optionsStack := stack.NewStack(NewOptionsHandler()).Use(r.RequestLogging, requestCounter, r.CORS)
	router.Handle("/user_preferences", optionsStack).Methods("OPTIONS").Name("OPTIONS /user_preferences")
	router.Handle("/user_preferences/{user_id}", optionsStack).Methods("OPTIONS").Name("OPTIONS /user_preferences/{user_id}")

	getPreferencesHandler := NewGetPreferencesHandler(r.PreferencesFinder, r.ErrorWriter)
	getPreferencesStack := stack.NewStack(getPreferencesHandler).Use(r.RequestLogging, requestCounter, r.CORS, r.NotificationPreferencesReadAuthenticator, r.DatabaseAllocator)
	router.Handle("/user_preferences", getPreferencesStack).Methods("GET").Name("GET /user_preferences")

	updatePreferencesHandler := NewUpdatePreferencesHandler(r.PreferenceUpdater, r.ErrorWriter)
	updatePreferencesStack := stack.NewStack(updatePreferencesHandler).Use(r.RequestLogging, requestCounter, r.CORS, r.NotificationPreferencesWriteAuthenticator, r.DatabaseAllocator)
	router.Handle("/user_preferences", updatePreferencesStack).Methods("PATCH").Name("PATCH /user_preferences")

	getUserPreferencesHandler := NewGetUserPreferencesHandler(r.PreferencesFinder, r.ErrorWriter)
	getUserPreferencesStack := stack.NewStack(getUserPreferencesHandler).Use(r.RequestLogging, requestCounter, r.CORS, r.NotificationPreferencesAdminAuthenticator, r.DatabaseAllocator)
	router.Handle("/user_preferences/{user_id}", getUserPreferencesStack).Methods("GET").Name("GET /user_preferences/{user_id}")

	updateUserPreferencesHandler := NewUpdateUserPreferencesHandler(r.PreferenceUpdater, r.ErrorWriter)
	updateUserPreferencesStack := stack.NewStack(updateUserPreferencesHandler).Use(r.RequestLogging, requestCounter, r.CORS, r.NotificationPreferencesAdminAuthenticator, r.DatabaseAllocator)
	router.Handle("/user_preferences/{user_id}", updateUserPreferencesStack).Methods("PATCH").Name("PATCH /user_preferences/{user_id}")
}
