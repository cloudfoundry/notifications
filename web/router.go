package web

import (
	"strings"

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
	PreferencesFinder() *services.PreferencesFinder
	PreferenceUpdater() services.PreferenceUpdater
	TemplateServiceObjects() (services.TemplateCreator, services.TemplateFinder, services.TemplateUpdater, services.TemplateDeleter, services.TemplateLister, services.TemplateAssigner)
	Database() models.DatabaseInterface
	Logging() stack.Middleware
	ErrorWriter() handlers.ErrorWriter
	Authenticator(...string) middleware.Authenticator
	CORS() middleware.CORS
}

type Router struct {
	stacks map[string]stack.Stack
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
	templateCreator, templateFinder, templateUpdater, templateDeleter, templateLister, templateAssigner := mother.TemplateServiceObjects()
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
	notificationsTemplateAdminAuthenticator := mother.Authenticator("notification_templates.admin")
	database := mother.Database()
	cors := mother.CORS()

	return Router{
		stacks: map[string]stack.Stack{
			"GET /info":                                                       stack.NewStack(handlers.NewGetInfo()).Use(logging),
			"POST /users/{guid}":                                              stack.NewStack(handlers.NewNotifyUser(notify, errorWriter, userStrategy, database)).Use(logging, notificationsWriteAuthenticator),
			"POST /spaces/{guid}":                                             stack.NewStack(handlers.NewNotifySpace(notify, errorWriter, spaceStrategy, database)).Use(logging, notificationsWriteAuthenticator),
			"POST /organizations/{guid}":                                      stack.NewStack(handlers.NewNotifyOrganization(notify, errorWriter, organizationStrategy, database)).Use(logging, notificationsWriteAuthenticator),
			"POST /everyone":                                                  stack.NewStack(handlers.NewNotifyEveryone(notify, errorWriter, everyoneStrategy, database)).Use(logging, notificationsWriteAuthenticator),
			"POST /uaa_scopes/{scope}":                                        stack.NewStack(handlers.NewNotifyUAAScope(notify, errorWriter, uaaScopeStrategy, database)).Use(logging, notificationsWriteAuthenticator),
			"POST /emails":                                                    stack.NewStack(handlers.NewNotifyEmail(notify, errorWriter, emailStrategy, database)).Use(logging, emailsWriteAuthenticator),
			"PUT /registration":                                               stack.NewStack(handlers.NewRegisterNotifications(registrar, errorWriter, database)).Use(logging, notificationsWriteAuthenticator),
			"PUT /notifications":                                              stack.NewStack(handlers.NewRegisterClientWithNotifications(registrar, errorWriter, database)).Use(logging, notificationsWriteAuthenticator),
			"GET /notifications":                                              stack.NewStack(handlers.NewGetAllNotifications(notificationsFinder, errorWriter)).Use(logging, notificationsManageAuthenticator),
			"OPTIONS /user_preferences":                                       stack.NewStack(handlers.NewOptionsPreferences()).Use(logging, cors),
			"OPTIONS /user_preferences/{guid}":                                stack.NewStack(handlers.NewOptionsPreferences()).Use(logging, cors),
			"GET /user_preferences":                                           stack.NewStack(handlers.NewGetPreferences(preferencesFinder, errorWriter)).Use(logging, cors, notificationPreferencesReadAuthenticator),
			"GET /user_preferences/{guid}":                                    stack.NewStack(handlers.NewGetPreferencesForUser(preferencesFinder, errorWriter)).Use(logging, cors, notificationPreferencesAdminAuthenticator),
			"PATCH /user_preferences":                                         stack.NewStack(handlers.NewUpdatePreferences(preferenceUpdater, errorWriter, database)).Use(logging, cors, notificationPreferencesWriteAuthenticator),
			"PATCH /user_preferences/{guid}":                                  stack.NewStack(handlers.NewUpdateSpecificUserPreferences(preferenceUpdater, errorWriter, database)).Use(logging, cors, notificationPreferencesAdminAuthenticator),
			"POST /templates":                                                 stack.NewStack(handlers.NewCreateTemplate(templateCreator, errorWriter)).Use(logging, notificationsTemplateWriteAuthenticator),
			"GET /templates/{templateID}":                                     stack.NewStack(handlers.NewGetTemplates(templateFinder, errorWriter)).Use(logging, notificationsTemplateReadAuthenticator),
			"PUT /templates/{templateID}":                                     stack.NewStack(handlers.NewUpdateTemplates(templateUpdater, errorWriter)).Use(logging, notificationsTemplateWriteAuthenticator),
			"PUT /deprecated_templates/{templateName}":                        stack.NewStack(handlers.NewSetTemplates(templateUpdater, errorWriter)).Use(logging, notificationsTemplateAdminAuthenticator),
			"DELETE /templates/{templateID}":                                  stack.NewStack(handlers.NewDeleteTemplates(templateDeleter, errorWriter)).Use(logging, notificationsTemplateWriteAuthenticator),
			"DELETE /deprecated_templates/{templateName}":                     stack.NewStack(handlers.NewUnsetTemplates(templateDeleter, errorWriter)).Use(logging, notificationsTemplateWriteAuthenticator),
			"GET /templates":                                                  stack.NewStack(handlers.NewListTemplates(templateLister, errorWriter)).Use(logging, notificationsTemplateReadAuthenticator),
			"PUT /clients/{clientID}/template":                                stack.NewStack(handlers.NewAssignClientTemplate(templateAssigner, errorWriter)).Use(logging, notificationsManageAuthenticator),
			"PUT /clients/{clientID}/notifications/{notificationID}/template": stack.NewStack(handlers.NewAssignNotificationTemplate(templateAssigner, errorWriter)).Use(logging, notificationsManageAuthenticator),
		},
	}
}

func (router Router) Routes() *mux.Router {
	r := mux.NewRouter()
	for methodPath, stack := range router.stacks {
		var name = methodPath
		parts := strings.SplitN(methodPath, " ", 2)
		r.Handle(parts[1], stack).Methods(parts[0]).Name(name)
	}
	return r
}
