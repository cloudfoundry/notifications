package web

import (
	"strings"

	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal/strategies"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/cloudfoundry-incubator/notifications/web/services"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

type MotherInterface interface {
	Registrar() services.Registrar
	EmailStrategy() strategies.EmailStrategy
	UserStrategy() strategies.UserStrategy
	SpaceStrategy() strategies.SpaceStrategy
	OrganizationStrategy() strategies.OrganizationStrategy
	EveryoneStrategy() strategies.EveryoneStrategy
	UAAScopeStrategy() strategies.UAAScopeStrategy
	NotificationsFinder() services.NotificationsFinder
	NotificationsUpdater() services.NotificationsUpdater
	PreferencesFinder() *services.PreferencesFinder
	PreferenceUpdater() services.PreferenceUpdater
	MessageFinder() services.MessageFinder
	TemplateServiceObjects() (services.TemplateCreator, services.TemplateFinder, services.TemplateUpdater, services.TemplateDeleter, services.TemplateLister, services.TemplateAssigner, services.TemplateAssociationLister)
	Database() models.DatabaseInterface
	Logging() middleware.RequestLogging
	ErrorWriter() handlers.ErrorWriter
	Authenticator(...string) middleware.Authenticator
	CORS() middleware.CORS
}

type Router struct {
	stacks map[string]stack.Stack
	router *mux.Router
}

func NewRouter(mother MotherInterface) Router {
	registrar := mother.Registrar()
	notificationsFinder := mother.NotificationsFinder()
	emailStrategy := mother.EmailStrategy()
	userStrategy := mother.UserStrategy()
	spaceStrategy := mother.SpaceStrategy()
	organizationStrategy := mother.OrganizationStrategy()
	everyoneStrategy := mother.EveryoneStrategy()
	uaaScopeStrategy := mother.UAAScopeStrategy()
	notify := handlers.NewNotify(mother.NotificationsFinder(), registrar)
	preferencesFinder := mother.PreferencesFinder()
	preferenceUpdater := mother.PreferenceUpdater()
	templateCreator, templateFinder, templateUpdater, templateDeleter, templateLister, templateAssigner, templateAssociationLister := mother.TemplateServiceObjects()
	notificationsUpdater := mother.NotificationsUpdater()
	messageFinder := mother.MessageFinder()
	logging := mother.Logging()
	errorWriter := mother.ErrorWriter()
	notificationsWriteAuthenticator := mother.Authenticator("notifications.write")
	notificationsManageAuthenticator := mother.Authenticator("notifications.manage")
	notificationPreferencesReadAuthenticator := mother.Authenticator("notification_preferences.read")
	notificationPreferencesWriteAuthenticator := mother.Authenticator("notification_preferences.write")
	notificationPreferencesAdminAuthenticator := mother.Authenticator("notification_preferences.admin")
	emailsWriteAuthenticator := mother.Authenticator("emails.write")
	notificationsTemplateWriteAuthenticator := mother.Authenticator("notification_templates.write")
	notificationsTemplateReadAuthenticator := mother.Authenticator("notification_templates.read")
	notificationsWriteOrEmailsWriteAuthenticator := mother.Authenticator("notifications.write", "emails.write")
	database := mother.Database()
	cors := mother.CORS()
	router := mux.NewRouter()
	requestCounter := middleware.NewRequestCounter(router, metrics.DefaultLogger)

	return Router{
		router: router,
		stacks: map[string]stack.Stack{
			"GET /info":                                                         stack.NewStack(handlers.NewGetInfo()).Use(logging, requestCounter),
			"POST /users/{user_id}":                                             stack.NewStack(handlers.NewNotifyUser(notify, errorWriter, userStrategy, database)).Use(logging, requestCounter, notificationsWriteAuthenticator),
			"POST /spaces/{space_id}":                                           stack.NewStack(handlers.NewNotifySpace(notify, errorWriter, spaceStrategy, database)).Use(logging, requestCounter, notificationsWriteAuthenticator),
			"POST /organizations/{org_id}":                                      stack.NewStack(handlers.NewNotifyOrganization(notify, errorWriter, organizationStrategy, database)).Use(logging, requestCounter, notificationsWriteAuthenticator),
			"POST /everyone":                                                    stack.NewStack(handlers.NewNotifyEveryone(notify, errorWriter, everyoneStrategy, database)).Use(logging, requestCounter, notificationsWriteAuthenticator),
			"POST /uaa_scopes/{scope}":                                          stack.NewStack(handlers.NewNotifyUAAScope(notify, errorWriter, uaaScopeStrategy, database)).Use(logging, requestCounter, notificationsWriteAuthenticator),
			"POST /emails":                                                      stack.NewStack(handlers.NewNotifyEmail(notify, errorWriter, emailStrategy, database)).Use(logging, requestCounter, emailsWriteAuthenticator),
			"PUT /registration":                                                 stack.NewStack(handlers.NewRegisterNotifications(registrar, errorWriter, database)).Use(logging, requestCounter, notificationsWriteAuthenticator),
			"PUT /notifications":                                                stack.NewStack(handlers.NewRegisterClientWithNotifications(registrar, errorWriter, database)).Use(logging, requestCounter, notificationsWriteAuthenticator),
			"GET /notifications":                                                stack.NewStack(handlers.NewGetAllNotifications(notificationsFinder, errorWriter)).Use(logging, requestCounter, notificationsManageAuthenticator),
			"OPTIONS /user_preferences":                                         stack.NewStack(handlers.NewOptionsPreferences()).Use(logging, requestCounter, cors),
			"OPTIONS /user_preferences/{user_id}":                               stack.NewStack(handlers.NewOptionsPreferences()).Use(logging, requestCounter, cors),
			"GET /user_preferences":                                             stack.NewStack(handlers.NewGetPreferences(preferencesFinder, errorWriter)).Use(logging, requestCounter, cors, notificationPreferencesReadAuthenticator),
			"GET /user_preferences/{user_id}":                                   stack.NewStack(handlers.NewGetPreferencesForUser(preferencesFinder, errorWriter)).Use(logging, requestCounter, cors, notificationPreferencesAdminAuthenticator),
			"PATCH /user_preferences":                                           stack.NewStack(handlers.NewUpdatePreferences(preferenceUpdater, errorWriter, database)).Use(logging, requestCounter, cors, notificationPreferencesWriteAuthenticator),
			"PATCH /user_preferences/{user_id}":                                 stack.NewStack(handlers.NewUpdateSpecificUserPreferences(preferenceUpdater, errorWriter, database)).Use(logging, requestCounter, cors, notificationPreferencesAdminAuthenticator),
			"GET /default_template":                                             stack.NewStack(handlers.NewGetDefaultTemplate(templateFinder, errorWriter)).Use(logging, requestCounter, notificationsTemplateReadAuthenticator),
			"PUT /default_template":                                             stack.NewStack(handlers.NewUpdateDefaultTemplate(templateUpdater, errorWriter)).Use(logging, requestCounter, notificationsTemplateWriteAuthenticator),
			"POST /templates":                                                   stack.NewStack(handlers.NewCreateTemplate(templateCreator, errorWriter)).Use(logging, requestCounter, notificationsTemplateWriteAuthenticator),
			"GET /templates/{template_id}":                                      stack.NewStack(handlers.NewGetTemplates(templateFinder, errorWriter)).Use(logging, requestCounter, notificationsTemplateReadAuthenticator),
			"PUT /templates/{template_id}":                                      stack.NewStack(handlers.NewUpdateTemplates(templateUpdater, errorWriter)).Use(logging, requestCounter, notificationsTemplateWriteAuthenticator),
			"DELETE /templates/{template_id}":                                   stack.NewStack(handlers.NewDeleteTemplates(templateDeleter, errorWriter)).Use(logging, requestCounter, notificationsTemplateWriteAuthenticator),
			"GET /templates/{template_id}/associations":                         stack.NewStack(handlers.NewListTemplateAssociations(templateAssociationLister, errorWriter)).Use(logging, requestCounter, notificationsManageAuthenticator),
			"GET /templates":                                                    stack.NewStack(handlers.NewListTemplates(templateLister, errorWriter)).Use(logging, requestCounter, notificationsTemplateReadAuthenticator),
			"PUT /clients/{client_id}/template":                                 stack.NewStack(handlers.NewAssignClientTemplate(templateAssigner, errorWriter)).Use(logging, requestCounter, notificationsManageAuthenticator),
			"PUT /clients/{client_id}/notifications/{notification_id}":          stack.NewStack(handlers.NewUpdateNotifications(notificationsUpdater, errorWriter)).Use(logging, requestCounter, notificationsManageAuthenticator),
			"PUT /clients/{client_id}/notifications/{notification_id}/template": stack.NewStack(handlers.NewAssignNotificationTemplate(templateAssigner, errorWriter)).Use(logging, requestCounter, notificationsManageAuthenticator),
			"GET /messages/{message_id}":                                        stack.NewStack(handlers.NewGetMessages(messageFinder, errorWriter)).Use(logging, requestCounter, notificationsWriteOrEmailsWriteAuthenticator),
		},
	}
}

func (router Router) Routes() *mux.Router {
	for methodPath, stack := range router.stacks {
		var name = methodPath
		parts := strings.SplitN(methodPath, " ", 2)
		router.router.Handle(parts[1], stack).Methods(parts[0]).Name(name)
	}
	return router.router
}
