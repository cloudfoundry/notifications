package web

import (
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

func NewClientsRouter(templateAssigner services.TemplateAssignerInterface,
	errorWriter handlers.ErrorWriterInterface,
	logging middleware.RequestLogging,
	notificationsManageAuthenticator middleware.Authenticator,
	databaseAllocator middleware.DatabaseAllocator,
	notificationsUpdater services.NotificationsUpdaterInterface) *mux.Router {

	router := mux.NewRouter()
	requestCounter := middleware.NewRequestCounter(router, metrics.DefaultLogger)

	router.Handle("/clients/{client_id}/template", stack.NewStack(handlers.NewAssignClientTemplate(templateAssigner, errorWriter)).Use(logging, requestCounter, notificationsManageAuthenticator, databaseAllocator)).Methods("PUT").Name("PUT /clients/{client_id}/template")
	router.Handle("/clients/{client_id}/notifications/{notification_id}", stack.NewStack(handlers.NewUpdateNotifications(notificationsUpdater, errorWriter)).Use(logging, requestCounter, notificationsManageAuthenticator, databaseAllocator)).Methods("PUT").Name("PUT /clients/{client_id}/notifications/{notification_id}")
	router.Handle("/clients/{client_id}/notifications/{notification_id}/template", stack.NewStack(handlers.NewAssignNotificationTemplate(templateAssigner, errorWriter)).Use(logging, requestCounter, notificationsManageAuthenticator, databaseAllocator)).Methods("PUT").Name("PUT /clients/{client_id}/notifications/{notification_id}/template")

	return router
}
