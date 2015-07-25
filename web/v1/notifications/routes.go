package notifications

import (
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

type Routes struct {
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

func (r Routes) Register(router *mux.Router) {
	requestCounter := middleware.NewRequestCounter(router, metrics.DefaultLogger)

	registrationHandler := NewRegistrationHandler(r.Registrar, r.ErrorWriter)
	registrationStack := stack.NewStack(registrationHandler).Use(r.RequestLogging, requestCounter, r.NotificationsWriteAuthenticator, r.DatabaseAllocator)
	router.Handle("/registration", registrationStack).Methods("PUT").Name("PUT /registration")

	putHandler := NewPutHandler(r.Registrar, r.ErrorWriter)
	putStack := stack.NewStack(putHandler).Use(r.RequestLogging, requestCounter, r.NotificationsWriteAuthenticator, r.DatabaseAllocator)
	router.Handle("/notifications", putStack).Methods("PUT").Name("PUT /notifications")

	listHandler := NewListHandler(r.NotificationsFinder, r.ErrorWriter)
	listStack := stack.NewStack(listHandler).Use(r.RequestLogging, requestCounter, r.NotificationsManageAuthenticator, r.DatabaseAllocator)
	router.Handle("/notifications", listStack).Methods("GET").Name("GET /notifications")

	updateHandler := NewUpdateHandler(r.NotificationsUpdater, r.ErrorWriter)
	updateStack := stack.NewStack(updateHandler).Use(r.RequestLogging, requestCounter, r.NotificationsManageAuthenticator, r.DatabaseAllocator)
	router.Handle("/clients/{client_id}/notifications/{notification_id}", updateStack).Methods("PUT").Name("PUT /clients/{client_id}/notifications/{notification_id}")

	assignTemplateHandler := NewAssignTemplateHandler(r.TemplateAssigner, r.ErrorWriter)
	assignTemplateStack := stack.NewStack(assignTemplateHandler).Use(r.RequestLogging, requestCounter, r.NotificationsManageAuthenticator, r.DatabaseAllocator)
	router.Handle("/clients/{client_id}/notifications/{notification_id}/template", assignTemplateStack).Methods("PUT").Name("PUT /clients/{client_id}/notifications/{notification_id}/template")
}
