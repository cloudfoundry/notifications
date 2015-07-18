package web

import (
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

type ClientsRouterConfig struct {
	RequestLogging                   RequestLogging
	NotificationsManageAuthenticator Authenticator
	DatabaseAllocator                DatabaseAllocator
	ErrorWriter                      handlers.ErrorWriterInterface

	TemplateAssigner     services.TemplateAssignerInterface
	NotificationsUpdater services.NotificationsUpdaterInterface
}

func NewClientsRouter(config ClientsRouterConfig) *mux.Router {
	router := mux.NewRouter()
	requestCounter := NewRequestCounter(router, metrics.DefaultLogger)

	assignClientTemplateHandler := handlers.NewAssignClientTemplate(config.TemplateAssigner, config.ErrorWriter)
	assignClientTemplateStack := stack.NewStack(assignClientTemplateHandler).Use(config.RequestLogging, requestCounter, config.NotificationsManageAuthenticator, config.DatabaseAllocator)
	router.Handle("/clients/{client_id}/template", assignClientTemplateStack).Methods("PUT").Name("PUT /clients/{client_id}/template")

	updateNotificationHandler := handlers.NewUpdateNotifications(config.NotificationsUpdater, config.ErrorWriter)
	updateNotificationStack := stack.NewStack(updateNotificationHandler).Use(config.RequestLogging, requestCounter, config.NotificationsManageAuthenticator, config.DatabaseAllocator)
	router.Handle("/clients/{client_id}/notifications/{notification_id}", updateNotificationStack).Methods("PUT").Name("PUT /clients/{client_id}/notifications/{notification_id}")

	assignNotificationTemplateHandler := handlers.NewAssignNotificationTemplate(config.TemplateAssigner, config.ErrorWriter)
	assignNotificationTemplateStack := stack.NewStack(assignNotificationTemplateHandler).Use(config.RequestLogging, requestCounter, config.NotificationsManageAuthenticator, config.DatabaseAllocator)
	router.Handle("/clients/{client_id}/notifications/{notification_id}/template", assignNotificationTemplateStack).Methods("PUT").Name("PUT /clients/{client_id}/notifications/{notification_id}/template")

	return router
}
