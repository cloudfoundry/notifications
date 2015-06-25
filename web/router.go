package web

import (
	"database/sql"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/postal/strategies"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
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

func NewRouter(mother MotherInterface, config Config) http.Handler {
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

	userPreferencesRouter := NewUserPreferencesRouter(logging, cors, preferencesFinder, errorWriter, notificationPreferencesReadAuthenticator, databaseAllocator, notificationPreferencesAdminAuthenticator, preferenceUpdater, notificationPreferencesWriteAuthenticator)
	clientsRouter := NewClientsRouter(templateAssigner, errorWriter, logging, notificationsManageAuthenticator, databaseAllocator, notificationsUpdater)
	messagesRouter := NewMessagesRouter(messageFinder, errorWriter, logging, notificationsWriteOrEmailsWriteAuthenticator, databaseAllocator)
	templatesRouter := NewTemplatesRouter(templateFinder, errorWriter, logging, notificationsTemplateReadAuthenticator, notificationsTemplateWriteAuthenticator, databaseAllocator, templateUpdater, templateCreator, templateDeleter, templateAssociationLister, notificationsManageAuthenticator, templateLister)
	notificationsRouter := NewNotificationsRouter(registrar, errorWriter, logging, notificationsWriteAuthenticator, databaseAllocator, notificationsFinder, notificationsManageAuthenticator)
	notifyRouter := NewNotifyRouter(notify, errorWriter, userStrategy, logging, notificationsWriteAuthenticator, databaseAllocator, spaceStrategy, organizationStrategy, everyoneStrategy, uaaScopeStrategy, emailStrategy, emailsWriteAuthenticator)

	pool := NewRouterPool()
	pool.AddMux(NewInfoRouter(logging))
	pool.AddMux(userPreferencesRouter)
	pool.AddMux(clientsRouter)
	pool.AddMux(messagesRouter)
	pool.AddMux(templatesRouter)
	pool.AddMux(notificationsRouter)
	pool.AddMux(notifyRouter)

	return pool
}
