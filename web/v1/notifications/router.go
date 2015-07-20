package notifications

import (
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

type RouterConfig struct {
	RequestLogging                   stack.Middleware
	DatabaseAllocator                stack.Middleware
	NotificationsWriteAuthenticator  stack.Middleware
	NotificationsManageAuthenticator stack.Middleware

	ErrorWriter          errorWriter
	Registrar            services.RegistrarInterface
	TemplateAssigner     services.TemplateAssignerInterface
	NotificationsFinder  services.NotificationsFinderInterface
	NotificationsUpdater services.NotificationsUpdaterInterface
}

func NewRouter(config RouterConfig) *mux.Router {
	router := mux.NewRouter()
	requestCounter := middleware.NewRequestCounter(router, metrics.DefaultLogger)

	registrationHandler := NewRegistrationHandler(config.Registrar, config.ErrorWriter)
	registrationStack := stack.NewStack(registrationHandler).Use(config.RequestLogging, requestCounter, config.NotificationsWriteAuthenticator, config.DatabaseAllocator)
	router.Handle("/registration", registrationStack).Methods("PUT").Name("PUT /registration")

	putHandler := NewPutHandler(config.Registrar, config.ErrorWriter)
	putStack := stack.NewStack(putHandler).Use(config.RequestLogging, requestCounter, config.NotificationsWriteAuthenticator, config.DatabaseAllocator)
	router.Handle("/notifications", putStack).Methods("PUT").Name("PUT /notifications")

	listHandler := NewListHandler(config.NotificationsFinder, config.ErrorWriter)
	listStack := stack.NewStack(listHandler).Use(config.RequestLogging, requestCounter, config.NotificationsManageAuthenticator, config.DatabaseAllocator)
	router.Handle("/notifications", listStack).Methods("GET").Name("GET /notifications")

	updateHandler := NewUpdateHandler(config.NotificationsUpdater, config.ErrorWriter)
	updateStack := stack.NewStack(updateHandler).Use(config.RequestLogging, requestCounter, config.NotificationsManageAuthenticator, config.DatabaseAllocator)
	router.Handle("/clients/{client_id}/notifications/{notification_id}", updateStack).Methods("PUT").Name("PUT /clients/{client_id}/notifications/{notification_id}")

	assignTemplateHandler := NewAssignTemplateHandler(config.TemplateAssigner, config.ErrorWriter)
	assignTemplateStack := stack.NewStack(assignTemplateHandler).Use(config.RequestLogging, requestCounter, config.NotificationsManageAuthenticator, config.DatabaseAllocator)
	router.Handle("/clients/{client_id}/notifications/{notification_id}/template", assignTemplateStack).Methods("PUT").Name("PUT /clients/{client_id}/notifications/{notification_id}/template")

	return router
}
