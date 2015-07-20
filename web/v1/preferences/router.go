package preferences

import (
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

type RouterConfig struct {
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

func NewRouter(config RouterConfig) *mux.Router {
	router := mux.NewRouter()
	requestCounter := middleware.NewRequestCounter(router, metrics.DefaultLogger)

	optionsStack := stack.NewStack(NewOptionsHandler()).Use(config.RequestLogging, requestCounter, config.CORS)
	router.Handle("/user_preferences", optionsStack).Methods("OPTIONS").Name("OPTIONS /user_preferences")
	router.Handle("/user_preferences/{user_id}", optionsStack).Methods("OPTIONS").Name("OPTIONS /user_preferences/{user_id}")

	getPreferencesHandler := NewGetPreferencesHandler(config.PreferencesFinder, config.ErrorWriter)
	getPreferencesStack := stack.NewStack(getPreferencesHandler).Use(config.RequestLogging, requestCounter, config.CORS, config.NotificationPreferencesReadAuthenticator, config.DatabaseAllocator)
	router.Handle("/user_preferences", getPreferencesStack).Methods("GET").Name("GET /user_preferences")

	updatePreferencesHandler := NewUpdatePreferencesHandler(config.PreferenceUpdater, config.ErrorWriter)
	updatePreferencesStack := stack.NewStack(updatePreferencesHandler).Use(config.RequestLogging, requestCounter, config.CORS, config.NotificationPreferencesWriteAuthenticator, config.DatabaseAllocator)
	router.Handle("/user_preferences", updatePreferencesStack).Methods("PATCH").Name("PATCH /user_preferences")

	getUserPreferencesHandler := NewGetUserPreferencesHandler(config.PreferencesFinder, config.ErrorWriter)
	getUserPreferencesStack := stack.NewStack(getUserPreferencesHandler).Use(config.RequestLogging, requestCounter, config.CORS, config.NotificationPreferencesAdminAuthenticator, config.DatabaseAllocator)
	router.Handle("/user_preferences/{user_id}", getUserPreferencesStack).Methods("GET").Name("GET /user_preferences/{user_id}")

	updateUserPreferencesHandler := NewUpdateUserPreferencesHandler(config.PreferenceUpdater, config.ErrorWriter)
	updateUserPreferencesStack := stack.NewStack(updateUserPreferencesHandler).Use(config.RequestLogging, requestCounter, config.CORS, config.NotificationPreferencesAdminAuthenticator, config.DatabaseAllocator)
	router.Handle("/user_preferences/{user_id}", updateUserPreferencesStack).Methods("PATCH").Name("PATCH /user_preferences/{user_id}")

	return router
}
