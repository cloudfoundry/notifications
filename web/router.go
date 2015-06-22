package web

import (
	"database/sql"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/metrics"
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
	Logging() middleware.RequestLogging
	ErrorWriter() handlers.ErrorWriter
	Authenticator(...string) middleware.Authenticator
	CORS() middleware.CORS
	SQLDatabase() *sql.DB
}

type Router struct {
	stacks map[string]stack.Stack
	router *mux.Router
}

func NewRouter(mother MotherInterface, config Config) Router {
	router := mux.NewRouter()

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
	databaseAllocator := middleware.NewDatabaseAllocator(mother.SQLDatabase(), config.DBLoggingEnabled)
	cors := mother.CORS()
	requestCounter := middleware.NewRequestCounter(router, metrics.DefaultLogger)

	userPreferencesRouter := NewUserPreferencesRouter(logging, cors, preferencesFinder, errorWriter, notificationPreferencesReadAuthenticator, databaseAllocator, notificationPreferencesAdminAuthenticator, preferenceUpdater, notificationPreferencesWriteAuthenticator)
	clientsRouter := NewClientsRouter(templateAssigner, errorWriter, logging, notificationsManageAuthenticator, databaseAllocator, notificationsUpdater)
	messagesRouter := NewMessagesRouter(messageFinder, errorWriter, logging, notificationsWriteOrEmailsWriteAuthenticator, databaseAllocator)
	templatesRouter := NewTemplatesRouter(templateFinder, errorWriter, logging, notificationsTemplateReadAuthenticator, notificationsTemplateWriteAuthenticator, databaseAllocator, templateUpdater, templateCreator, templateDeleter, templateAssociationLister, notificationsManageAuthenticator, templateLister)

	router.Handle("/info", NewInfoRouter(logging))
	router.Handle("/user_preferences{anything:.*}", userPreferencesRouter)
	router.Handle("/clients{anything:.*}", clientsRouter)
	router.Handle("/messages{anything:.*}", messagesRouter)
	router.Handle("/default_template{anything:.*}", templatesRouter)
	router.Handle("/templates{anything:.*}", templatesRouter)

	registrationStack := newRegistrationStack(registrar, errorWriter, logging, requestCounter, notificationsWriteAuthenticator, databaseAllocator, notificationsFinder, notificationsManageAuthenticator)
	notificationsStack := newNotificationsStack(notify, errorWriter, userStrategy, logging, requestCounter, notificationsWriteAuthenticator, databaseAllocator, spaceStrategy, organizationStrategy, everyoneStrategy, uaaScopeStrategy, emailStrategy, emailsWriteAuthenticator)

	stacks := make(map[string]stack.Stack)
	for _, s := range []map[string]stack.Stack{notificationsStack, registrationStack} {
		for route, handler := range s {
			stacks[route] = handler
		}
	}

	return Router{
		router: router,
		stacks: stacks,
	}
}

func newNotificationsStack(notify handlers.Notify, errorWriter handlers.ErrorWriter, userStrategy strategies.UserStrategy, logging middleware.RequestLogging, requestCounter middleware.RequestCounter, notificationsWriteAuthenticator middleware.Authenticator, databaseAllocator middleware.DatabaseAllocator, spaceStrategy strategies.SpaceStrategy, organizationStrategy strategies.OrganizationStrategy, everyoneStrategy strategies.EveryoneStrategy, uaaScopeStrategy strategies.UAAScopeStrategy, emailStrategy strategies.EmailStrategy, emailsWriteAuthenticator middleware.Authenticator) map[string]stack.Stack {
	return map[string]stack.Stack{
		"POST /users/{user_id}":        stack.NewStack(handlers.NewNotifyUser(notify, errorWriter, userStrategy)).Use(logging, requestCounter, notificationsWriteAuthenticator, databaseAllocator),
		"POST /spaces/{space_id}":      stack.NewStack(handlers.NewNotifySpace(notify, errorWriter, spaceStrategy)).Use(logging, requestCounter, notificationsWriteAuthenticator, databaseAllocator),
		"POST /organizations/{org_id}": stack.NewStack(handlers.NewNotifyOrganization(notify, errorWriter, organizationStrategy)).Use(logging, requestCounter, notificationsWriteAuthenticator, databaseAllocator),
		"POST /everyone":               stack.NewStack(handlers.NewNotifyEveryone(notify, errorWriter, everyoneStrategy)).Use(logging, requestCounter, notificationsWriteAuthenticator, databaseAllocator),
		"POST /uaa_scopes/{scope}":     stack.NewStack(handlers.NewNotifyUAAScope(notify, errorWriter, uaaScopeStrategy)).Use(logging, requestCounter, notificationsWriteAuthenticator, databaseAllocator),
		"POST /emails":                 stack.NewStack(handlers.NewNotifyEmail(notify, errorWriter, emailStrategy)).Use(logging, requestCounter, emailsWriteAuthenticator, databaseAllocator),
	}
}

func newRegistrationStack(registrar services.Registrar, errorWriter handlers.ErrorWriter, logging middleware.RequestLogging, requestCounter middleware.RequestCounter, notificationsWriteAuthenticator middleware.Authenticator, databaseAllocator middleware.DatabaseAllocator, notificationsFinder services.NotificationsFinderInterface, notificationsManageAuthenticator middleware.Authenticator) map[string]stack.Stack {
	return map[string]stack.Stack{
		"PUT /registration":  stack.NewStack(handlers.NewRegisterNotifications(registrar, errorWriter)).Use(logging, requestCounter, notificationsWriteAuthenticator, databaseAllocator),
		"PUT /notifications": stack.NewStack(handlers.NewRegisterClientWithNotifications(registrar, errorWriter)).Use(logging, requestCounter, notificationsWriteAuthenticator, databaseAllocator),
		"GET /notifications": stack.NewStack(handlers.NewGetAllNotifications(notificationsFinder, errorWriter)).Use(logging, requestCounter, notificationsManageAuthenticator, databaseAllocator),
	}
}

func (r Router) Routes() *mux.Router {
	for methodPath, stack := range r.stacks {
		var name = methodPath
		parts := strings.SplitN(methodPath, " ", 2)
		r.router.Handle(parts[1], stack).Methods(parts[0]).Name(name)
	}

	return r.router
}
