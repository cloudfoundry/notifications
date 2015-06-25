package web

import (
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

func NewUserPreferencesRouter(logging middleware.RequestLogging,
	cors middleware.CORS,
	preferencesFinder services.PreferencesFinderInterface,
	errorWriter handlers.ErrorWriterInterface,
	notificationPreferencesReadAuthenticator middleware.Authenticator,
	databaseAllocator middleware.DatabaseAllocator,
	notificationPreferencesAdminAuthenticator middleware.Authenticator,
	preferenceUpdater services.PreferenceUpdaterInterface,
	notificationPreferencesWriteAuthenticator middleware.Authenticator) *mux.Router {

	router := mux.NewRouter()
	requestCounter := middleware.NewRequestCounter(router, metrics.DefaultLogger)

	router.Handle("/user_preferences", stack.NewStack(handlers.NewOptionsPreferences()).Use(logging, requestCounter, cors)).Methods("OPTIONS").Name("OPTIONS /user_preferences")
	router.Handle("/user_preferences", stack.NewStack(handlers.NewGetPreferences(preferencesFinder, errorWriter)).Use(logging, requestCounter, cors, notificationPreferencesReadAuthenticator, databaseAllocator)).Methods("GET").Name("GET /user_preferences")
	router.Handle("/user_preferences", stack.NewStack(handlers.NewUpdatePreferences(preferenceUpdater, errorWriter)).Use(logging, requestCounter, cors, notificationPreferencesWriteAuthenticator, databaseAllocator)).Methods("PATCH").Name("PATCH /user_preferences")

	router.Handle("/user_preferences/{user_id}", stack.NewStack(handlers.NewOptionsPreferences()).Use(logging, requestCounter, cors)).Methods("OPTIONS").Name("OPTIONS /user_preferences/{user_id}")
	router.Handle("/user_preferences/{user_id}", stack.NewStack(handlers.NewGetPreferencesForUser(preferencesFinder, errorWriter)).Use(logging, requestCounter, cors, notificationPreferencesAdminAuthenticator, databaseAllocator)).Methods("GET").Name("GET /user_preferences/{user_id}")
	router.Handle("/user_preferences/{user_id}", stack.NewStack(handlers.NewUpdateSpecificUserPreferences(preferenceUpdater, errorWriter)).Use(logging, requestCounter, cors, notificationPreferencesAdminAuthenticator, databaseAllocator)).Methods("PATCH").Name("PATCH /user_preferences/{user_id}")

	return router
}
