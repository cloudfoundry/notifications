package web

import (
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

type UserPreferencesRouterConfig struct {
	ErrorWriter       handlers.ErrorWriterInterface
	PreferencesFinder services.PreferencesFinderInterface
	PreferenceUpdater services.PreferenceUpdaterInterface

	CORS                                      CORS
	RequestLogging                            RequestLogging
	DatabaseAllocator                         DatabaseAllocator
	NotificationPreferencesReadAuthenticator  Authenticator
	NotificationPreferencesAdminAuthenticator Authenticator
	NotificationPreferencesWriteAuthenticator Authenticator
}

func NewUserPreferencesRouter(config UserPreferencesRouterConfig) *mux.Router {
	router := mux.NewRouter()
	requestCounter := NewRequestCounter(router, metrics.DefaultLogger)

	optionsPreferencesStack := stack.NewStack(handlers.NewOptionsPreferences()).Use(config.RequestLogging, requestCounter, config.CORS)
	router.Handle("/user_preferences", optionsPreferencesStack).Methods("OPTIONS").Name("OPTIONS /user_preferences")
	router.Handle("/user_preferences/{user_id}", optionsPreferencesStack).Methods("OPTIONS").Name("OPTIONS /user_preferences/{user_id}")

	getPreferencesHandler := handlers.NewGetPreferences(config.PreferencesFinder, config.ErrorWriter)
	getPreferencesStack := stack.NewStack(getPreferencesHandler).Use(config.RequestLogging, requestCounter, config.CORS, config.NotificationPreferencesReadAuthenticator, config.DatabaseAllocator)
	router.Handle("/user_preferences", getPreferencesStack).Methods("GET").Name("GET /user_preferences")

	updatePreferencesHandler := handlers.NewUpdatePreferences(config.PreferenceUpdater, config.ErrorWriter)
	updatePreferencesStack := stack.NewStack(updatePreferencesHandler).Use(config.RequestLogging, requestCounter, config.CORS, config.NotificationPreferencesWriteAuthenticator, config.DatabaseAllocator)
	router.Handle("/user_preferences", updatePreferencesStack).Methods("PATCH").Name("PATCH /user_preferences")

	getPreferencesForUserHandler := handlers.NewGetPreferencesForUser(config.PreferencesFinder, config.ErrorWriter)
	getPreferencesForUserStack := stack.NewStack(getPreferencesForUserHandler).Use(config.RequestLogging, requestCounter, config.CORS, config.NotificationPreferencesAdminAuthenticator, config.DatabaseAllocator)
	router.Handle("/user_preferences/{user_id}", getPreferencesForUserStack).Methods("GET").Name("GET /user_preferences/{user_id}")

	updateSpecificUserPreferencesHandler := handlers.NewUpdateSpecificUserPreferences(config.PreferenceUpdater, config.ErrorWriter)
	updateSpecificUserPreferencesStack := stack.NewStack(updateSpecificUserPreferencesHandler).Use(config.RequestLogging, requestCounter, config.CORS, config.NotificationPreferencesAdminAuthenticator, config.DatabaseAllocator)
	router.Handle("/user_preferences/{user_id}", updateSpecificUserPreferencesStack).Methods("PATCH").Name("PATCH /user_preferences/{user_id}")

	return router
}
