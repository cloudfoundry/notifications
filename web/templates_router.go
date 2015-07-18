package web

import (
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

type TemplatesRouterConfig struct {
	TemplateFinder            services.TemplateFinderInterface
	TemplateLister            services.TemplateListerInterface
	TemplateUpdater           services.TemplateUpdaterInterface
	TemplateCreator           services.TemplateCreatorInterface
	TemplateDeleter           services.TemplateDeleterInterface
	TemplateAssociationLister services.TemplateAssociationListerInterface
	ErrorWriter               handlers.ErrorWriterInterface

	RequestLogging                          RequestLogging
	DatabaseAllocator                       DatabaseAllocator
	NotificationTemplatesReadAuthenticator  Authenticator
	NotificationTemplatesWriteAuthenticator Authenticator
	NotificationsManageAuthenticator        Authenticator
}

func NewTemplatesRouter(config TemplatesRouterConfig) *mux.Router {
	router := mux.NewRouter()
	requestCounter := NewRequestCounter(router, metrics.DefaultLogger)

	getDefaultTemplateHandler := handlers.NewGetDefaultTemplate(config.TemplateFinder, config.ErrorWriter)
	getDefaultTemplateStack := stack.NewStack(getDefaultTemplateHandler).Use(config.RequestLogging, requestCounter, config.NotificationTemplatesReadAuthenticator, config.DatabaseAllocator)
	router.Handle("/default_template", getDefaultTemplateStack).Methods("GET").Name("GET /default_template")

	updateDefaultTemplateHandler := handlers.NewUpdateDefaultTemplate(config.TemplateUpdater, config.ErrorWriter)
	updateDefaultTemplateStack := stack.NewStack(updateDefaultTemplateHandler).Use(config.RequestLogging, requestCounter, config.NotificationTemplatesWriteAuthenticator, config.DatabaseAllocator)
	router.Handle("/default_template", updateDefaultTemplateStack).Methods("PUT").Name("PUT /default_template")

	listTemplatesHandler := handlers.NewListTemplates(config.TemplateLister, config.ErrorWriter)
	listTemplatesStack := stack.NewStack(listTemplatesHandler).Use(config.RequestLogging, requestCounter, config.NotificationTemplatesReadAuthenticator, config.DatabaseAllocator)
	router.Handle("/templates", listTemplatesStack).Methods("GET").Name("GET /templates")

	createTemplateHandler := handlers.NewCreateTemplate(config.TemplateCreator, config.ErrorWriter)
	createTemplateStack := stack.NewStack(createTemplateHandler).Use(config.RequestLogging, requestCounter, config.NotificationTemplatesWriteAuthenticator, config.DatabaseAllocator)
	router.Handle("/templates", createTemplateStack).Methods("POST").Name("POST /templates")

	getTemplateHandler := handlers.NewGetTemplates(config.TemplateFinder, config.ErrorWriter)
	getTemplateStack := stack.NewStack(getTemplateHandler).Use(config.RequestLogging, requestCounter, config.NotificationTemplatesReadAuthenticator, config.DatabaseAllocator)
	router.Handle("/templates/{template_id}", getTemplateStack).Methods("GET").Name("GET /templates/{template_id}")

	updateTemplateHandler := handlers.NewUpdateTemplates(config.TemplateUpdater, config.ErrorWriter)
	updateTemplateStack := stack.NewStack(updateTemplateHandler).Use(config.RequestLogging, requestCounter, config.NotificationTemplatesWriteAuthenticator, config.DatabaseAllocator)
	router.Handle("/templates/{template_id}", updateTemplateStack).Methods("PUT").Name("PUT /templates/{template_id}")

	deleteTemplateHandler := handlers.NewDeleteTemplates(config.TemplateDeleter, config.ErrorWriter)
	deleteTemplateStack := stack.NewStack(deleteTemplateHandler).Use(config.RequestLogging, requestCounter, config.NotificationTemplatesWriteAuthenticator, config.DatabaseAllocator)
	router.Handle("/templates/{template_id}", deleteTemplateStack).Methods("DELETE").Name("DELETE /templates/{template_id}")

	listTemplateAssociationsHandler := handlers.NewListTemplateAssociations(config.TemplateAssociationLister, config.ErrorWriter)
	listTemplateAssociationsStack := stack.NewStack(listTemplateAssociationsHandler).Use(config.RequestLogging, requestCounter, config.NotificationsManageAuthenticator, config.DatabaseAllocator)
	router.Handle("/templates/{template_id}/associations", listTemplateAssociationsStack).Methods("GET").Name("GET /templates/{template_id}/associations")

	return router
}
