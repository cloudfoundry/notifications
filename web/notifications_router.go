package web

import (
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

func NewNotificationsRouter(registrar services.RegistrarInterface,
	errorWriter handlers.ErrorWriterInterface,
	logging middleware.RequestLogging,
	notificationsWriteAuthenticator middleware.Authenticator,
	databaseAllocator middleware.DatabaseAllocator,
	notificationsFinder services.NotificationsFinderInterface,
	notificationsManageAuthenticator middleware.Authenticator) *mux.Router {

	router := mux.NewRouter()
	requestCounter := middleware.NewRequestCounter(router, metrics.DefaultLogger)

	router.Handle("/registration", stack.NewStack(handlers.NewRegisterNotifications(registrar, errorWriter)).Use(logging, requestCounter, notificationsWriteAuthenticator, databaseAllocator)).Methods("PUT").Name("PUT /registration")
	router.Handle("/notifications", stack.NewStack(handlers.NewRegisterClientWithNotifications(registrar, errorWriter)).Use(logging, requestCounter, notificationsWriteAuthenticator, databaseAllocator)).Methods("PUT").Name("PUT /notifications")
	router.Handle("/notifications", stack.NewStack(handlers.NewGetAllNotifications(notificationsFinder, errorWriter)).Use(logging, requestCounter, notificationsManageAuthenticator, databaseAllocator)).Methods("GET").Name("GET /notifications")

	return router
}
