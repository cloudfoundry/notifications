package templates

import (
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

type Routes struct {
	RequestLogging                          stack.Middleware
	DatabaseAllocator                       stack.Middleware
	NotificationTemplatesReadAuthenticator  stack.Middleware
	NotificationTemplatesWriteAuthenticator stack.Middleware
	NotificationsManageAuthenticator        stack.Middleware

	ErrorWriter               errorWriter
	TemplateFinder            services.TemplateFinderInterface
	TemplateLister            services.TemplateListerInterface
	TemplateUpdater           services.TemplateUpdaterInterface
	TemplateCreator           services.TemplateCreatorInterface
	TemplateDeleter           services.TemplateDeleterInterface
	TemplateAssociationLister services.TemplateAssociationListerInterface
}

func (r Routes) Register(router *mux.Router) {
	requestCounter := middleware.NewRequestCounter(router, metrics.DefaultLogger)

	getDefaultHandler := NewGetDefaultHandler(r.TemplateFinder, r.ErrorWriter)
	getDefaultStack := stack.NewStack(getDefaultHandler).Use(r.RequestLogging, requestCounter, r.NotificationTemplatesReadAuthenticator, r.DatabaseAllocator)
	router.Handle("/default_template", getDefaultStack).Methods("GET").Name("GET /default_template")

	updateDefaultHandler := NewUpdateDefaultHandler(r.TemplateUpdater, r.ErrorWriter)
	updateDefaultStack := stack.NewStack(updateDefaultHandler).Use(r.RequestLogging, requestCounter, r.NotificationTemplatesWriteAuthenticator, r.DatabaseAllocator)
	router.Handle("/default_template", updateDefaultStack).Methods("PUT").Name("PUT /default_template")

	listHandler := NewListHandler(r.TemplateLister, r.ErrorWriter)
	listStack := stack.NewStack(listHandler).Use(r.RequestLogging, requestCounter, r.NotificationTemplatesReadAuthenticator, r.DatabaseAllocator)
	router.Handle("/templates", listStack).Methods("GET").Name("GET /templates")

	createHandler := NewCreateHandler(r.TemplateCreator, r.ErrorWriter)
	createStack := stack.NewStack(createHandler).Use(r.RequestLogging, requestCounter, r.NotificationTemplatesWriteAuthenticator, r.DatabaseAllocator)
	router.Handle("/templates", createStack).Methods("POST").Name("POST /templates")

	getHandler := NewGetHandler(r.TemplateFinder, r.ErrorWriter)
	getStack := stack.NewStack(getHandler).Use(r.RequestLogging, requestCounter, r.NotificationTemplatesReadAuthenticator, r.DatabaseAllocator)
	router.Handle("/templates/{template_id}", getStack).Methods("GET").Name("GET /templates/{template_id}")

	updateHandler := NewUpdateHandler(r.TemplateUpdater, r.ErrorWriter)
	updateStack := stack.NewStack(updateHandler).Use(r.RequestLogging, requestCounter, r.NotificationTemplatesWriteAuthenticator, r.DatabaseAllocator)
	router.Handle("/templates/{template_id}", updateStack).Methods("PUT").Name("PUT /templates/{template_id}")

	deleteHandler := NewDeleteHandler(r.TemplateDeleter, r.ErrorWriter)
	deleteStack := stack.NewStack(deleteHandler).Use(r.RequestLogging, requestCounter, r.NotificationTemplatesWriteAuthenticator, r.DatabaseAllocator)
	router.Handle("/templates/{template_id}", deleteStack).Methods("DELETE").Name("DELETE /templates/{template_id}")

	listAssociationsHandler := NewListAssociationsHandler(r.TemplateAssociationLister, r.ErrorWriter)
	listAssociationsStack := stack.NewStack(listAssociationsHandler).Use(r.RequestLogging, requestCounter, r.NotificationsManageAuthenticator, r.DatabaseAllocator)
	router.Handle("/templates/{template_id}/associations", listAssociationsStack).Methods("GET").Name("GET /templates/{template_id}/associations")
}
