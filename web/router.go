package web

import (
	"database/sql"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/nu7hatch/gouuid"
)

type MotherInterface interface {
	Registrar() services.Registrar
	EmailStrategy() services.EmailStrategy
	UserStrategy() services.UserStrategy
	SpaceStrategy() services.SpaceStrategy
	OrganizationStrategy() services.OrganizationStrategy
	EveryoneStrategy() services.EveryoneStrategy
	UAAScopeStrategy() services.UAAScopeStrategy
	NotificationsFinder() services.NotificationsFinder
	NotificationsUpdater() services.NotificationsUpdater
	PreferencesFinder() *services.PreferencesFinder
	PreferenceUpdater() services.PreferenceUpdater
	MessageFinder() services.MessageFinder
	TemplateServiceObjects() (services.TemplateCreator, services.TemplateFinder, services.TemplateUpdater, services.TemplateDeleter, services.TemplateLister, services.TemplateAssigner, services.TemplateAssociationLister)
	Logging() RequestLogging
	ErrorWriter() handlers.ErrorWriter
	Authenticator(...string) Authenticator
	CORS() CORS
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
	databaseAllocator := NewDatabaseAllocator(mother.SQLDatabase(), config.DBLoggingEnabled)
	cors := mother.CORS()

	userPreferencesRouter := NewUserPreferencesRouter(logging, cors, preferencesFinder, errorWriter, notificationPreferencesReadAuthenticator, databaseAllocator, notificationPreferencesAdminAuthenticator, preferenceUpdater, notificationPreferencesWriteAuthenticator)
	clientsRouter := NewClientsRouter(templateAssigner, errorWriter, logging, notificationsManageAuthenticator, databaseAllocator, notificationsUpdater)
	messagesRouter := NewMessagesRouter(messageFinder, errorWriter, logging, notificationsWriteOrEmailsWriteAuthenticator, databaseAllocator)
	templatesRouter := NewTemplatesRouter(templateFinder, errorWriter, logging, notificationsTemplateReadAuthenticator, notificationsTemplateWriteAuthenticator, databaseAllocator, templateUpdater, templateCreator, templateDeleter, templateAssociationLister, notificationsManageAuthenticator, templateLister)
	notificationsRouter := NewNotificationsRouter(registrar, errorWriter, logging, notificationsWriteAuthenticator, databaseAllocator, notificationsFinder, notificationsManageAuthenticator)
	notifyRouter := NewNotifyRouter(notify, errorWriter, userStrategy, logging, notificationsWriteAuthenticator, databaseAllocator, spaceStrategy, organizationStrategy, everyoneStrategy, uaaScopeStrategy, emailStrategy, emailsWriteAuthenticator)

	v1 := NewRouterPool()
	v1.AddMux(NewInfoRouter(1, logging))
	v1.AddMux(userPreferencesRouter)
	v1.AddMux(clientsRouter)
	v1.AddMux(messagesRouter)
	v1.AddMux(templatesRouter)
	v1.AddMux(notificationsRouter)
	v1.AddMux(notifyRouter)

	sendersCollection := collections.NewSendersCollection(models.NewSendersRepository(uuid.NewV4))

	v2 := NewRouterPool()
	v2.AddMux(NewInfoRouter(2, logging))
	v2.AddMux(NewSendersRouter(logging, notificationsWriteAuthenticator, databaseAllocator, sendersCollection))
	v2.AddMux(NewNotificationTypesRouter(collections.NotificationTypesCollection{}))

	return VersionRouter{
		1: v1,
		2: v2,
	}
}
