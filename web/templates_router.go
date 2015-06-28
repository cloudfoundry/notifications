package web

import (
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

func NewTemplatesRouter(templateFinder services.TemplateFinderInterface,
	errorWriter handlers.ErrorWriterInterface,
	logging RequestLogging,
	notificationsTemplateReadAuthenticator Authenticator,
	notificationsTemplateWriteAuthenticator Authenticator,
	databaseAllocator DatabaseAllocator,
	templateUpdater services.TemplateUpdaterInterface,
	templateCreator services.TemplateCreatorInterface,
	templateDeleter services.TemplateDeleterInterface,
	templateAssociationLister services.TemplateAssociationListerInterface,
	notificationsManageAuthenticator Authenticator,
	templateLister services.TemplateListerInterface) *mux.Router {

	router := mux.NewRouter()
	requestCounter := NewRequestCounter(router, metrics.DefaultLogger)

	router.Handle("/default_template", stack.NewStack(handlers.NewGetDefaultTemplate(templateFinder, errorWriter)).Use(logging, requestCounter, notificationsTemplateReadAuthenticator, databaseAllocator)).Methods("GET").Name("GET /default_template")
	router.Handle("/default_template", stack.NewStack(handlers.NewUpdateDefaultTemplate(templateUpdater, errorWriter)).Use(logging, requestCounter, notificationsTemplateWriteAuthenticator, databaseAllocator)).Methods("PUT").Name("PUT /default_template")

	router.Handle("/templates", stack.NewStack(handlers.NewListTemplates(templateLister, errorWriter)).Use(logging, requestCounter, notificationsTemplateReadAuthenticator, databaseAllocator)).Methods("GET").Name("GET /templates")
	router.Handle("/templates", stack.NewStack(handlers.NewCreateTemplate(templateCreator, errorWriter)).Use(logging, requestCounter, notificationsTemplateWriteAuthenticator, databaseAllocator)).Methods("POST").Name("POST /templates")

	router.Handle("/templates/{template_id}", stack.NewStack(handlers.NewGetTemplates(templateFinder, errorWriter)).Use(logging, requestCounter, notificationsTemplateReadAuthenticator, databaseAllocator)).Methods("GET").Name("GET /templates/{template_id}")
	router.Handle("/templates/{template_id}", stack.NewStack(handlers.NewUpdateTemplates(templateUpdater, errorWriter)).Use(logging, requestCounter, notificationsTemplateWriteAuthenticator, databaseAllocator)).Methods("PUT").Name("PUT /templates/{template_id}")
	router.Handle("/templates/{template_id}", stack.NewStack(handlers.NewDeleteTemplates(templateDeleter, errorWriter)).Use(logging, requestCounter, notificationsTemplateWriteAuthenticator, databaseAllocator)).Methods("DELETE").Name("DELETE /templates/{template_id}")
	router.Handle("/templates/{template_id}/associations", stack.NewStack(handlers.NewListTemplateAssociations(templateAssociationLister, errorWriter)).Use(logging, requestCounter, notificationsManageAuthenticator, databaseAllocator)).Methods("GET").Name("GET /templates/{template_id}/associations")

	return router
}
