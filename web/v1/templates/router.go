package templates

import (
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

type RouterConfig struct {
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

func NewRouter(config RouterConfig) *mux.Router {
	router := mux.NewRouter()
	requestCounter := middleware.NewRequestCounter(router, metrics.DefaultLogger)

	getDefaultHandler := NewGetDefaultHandler(config.TemplateFinder, config.ErrorWriter)
	getDefaultStack := stack.NewStack(getDefaultHandler).Use(config.RequestLogging, requestCounter, config.NotificationTemplatesReadAuthenticator, config.DatabaseAllocator)
	router.Handle("/default_template", getDefaultStack).Methods("GET").Name("GET /default_template")

	updateDefaultHandler := NewUpdateDefaultHandler(config.TemplateUpdater, config.ErrorWriter)
	updateDefaultStack := stack.NewStack(updateDefaultHandler).Use(config.RequestLogging, requestCounter, config.NotificationTemplatesWriteAuthenticator, config.DatabaseAllocator)
	router.Handle("/default_template", updateDefaultStack).Methods("PUT").Name("PUT /default_template")

	listHandler := NewListHandler(config.TemplateLister, config.ErrorWriter)
	listStack := stack.NewStack(listHandler).Use(config.RequestLogging, requestCounter, config.NotificationTemplatesReadAuthenticator, config.DatabaseAllocator)
	router.Handle("/templates", listStack).Methods("GET").Name("GET /templates")

	createHandler := NewCreateHandler(config.TemplateCreator, config.ErrorWriter)
	createStack := stack.NewStack(createHandler).Use(config.RequestLogging, requestCounter, config.NotificationTemplatesWriteAuthenticator, config.DatabaseAllocator)
	router.Handle("/templates", createStack).Methods("POST").Name("POST /templates")

	getHandler := NewGetHandler(config.TemplateFinder, config.ErrorWriter)
	getStack := stack.NewStack(getHandler).Use(config.RequestLogging, requestCounter, config.NotificationTemplatesReadAuthenticator, config.DatabaseAllocator)
	router.Handle("/templates/{template_id}", getStack).Methods("GET").Name("GET /templates/{template_id}")

	updateHandler := NewUpdateHandler(config.TemplateUpdater, config.ErrorWriter)
	updateStack := stack.NewStack(updateHandler).Use(config.RequestLogging, requestCounter, config.NotificationTemplatesWriteAuthenticator, config.DatabaseAllocator)
	router.Handle("/templates/{template_id}", updateStack).Methods("PUT").Name("PUT /templates/{template_id}")

	deleteHandler := NewDeleteHandler(config.TemplateDeleter, config.ErrorWriter)
	deleteStack := stack.NewStack(deleteHandler).Use(config.RequestLogging, requestCounter, config.NotificationTemplatesWriteAuthenticator, config.DatabaseAllocator)
	router.Handle("/templates/{template_id}", deleteStack).Methods("DELETE").Name("DELETE /templates/{template_id}")

	listAssociationsHandler := NewListAssociationsHandler(config.TemplateAssociationLister, config.ErrorWriter)
	listAssociationsStack := stack.NewStack(listAssociationsHandler).Use(config.RequestLogging, requestCounter, config.NotificationsManageAuthenticator, config.DatabaseAllocator)
	router.Handle("/templates/{template_id}/associations", listAssociationsStack).Methods("GET").Name("GET /templates/{template_id}/associations")

	return router
}
