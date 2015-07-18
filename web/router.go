package web

import (
	"database/sql"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/cloudfoundry-incubator/notifications/web/v2/notificationtypes"
	"github.com/cloudfoundry-incubator/notifications/web/v2/senders"
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

	v1 := NewRouterPool()
	v1.AddMux(NewInfoRouter(InfoRouterConfig{
		Version:        1,
		RequestLogging: logging,
	}))
	v1.AddMux(NewUserPreferencesRouter(UserPreferencesRouterConfig{
		ErrorWriter:       errorWriter,
		PreferencesFinder: preferencesFinder,
		PreferenceUpdater: preferenceUpdater,

		CORS:                                      cors,
		RequestLogging:                            logging,
		DatabaseAllocator:                         databaseAllocator,
		NotificationPreferencesReadAuthenticator:  notificationPreferencesReadAuthenticator,
		NotificationPreferencesWriteAuthenticator: notificationPreferencesWriteAuthenticator,
		NotificationPreferencesAdminAuthenticator: notificationPreferencesAdminAuthenticator,
	}))
	v1.AddMux(NewClientsRouter(ClientsRouterConfig{
		ErrorWriter:          errorWriter,
		TemplateAssigner:     templateAssigner,
		NotificationsUpdater: notificationsUpdater,

		RequestLogging:                   logging,
		DatabaseAllocator:                databaseAllocator,
		NotificationsManageAuthenticator: notificationsManageAuthenticator,
	}))
	v1.AddMux(NewMessagesRouter(MessagesRouterConfig{
		ErrorWriter:   errorWriter,
		MessageFinder: messageFinder,

		RequestLogging:                               logging,
		DatabaseAllocator:                            databaseAllocator,
		NotificationsWriteOrEmailsWriteAuthenticator: notificationsWriteOrEmailsWriteAuthenticator,
	}))
	v1.AddMux(NewTemplatesRouter(TemplatesRouterConfig{
		ErrorWriter:               errorWriter,
		TemplateFinder:            templateFinder,
		TemplateUpdater:           templateUpdater,
		TemplateCreator:           templateCreator,
		TemplateDeleter:           templateDeleter,
		TemplateLister:            templateLister,
		TemplateAssociationLister: templateAssociationLister,

		RequestLogging:                          logging,
		DatabaseAllocator:                       databaseAllocator,
		NotificationTemplatesReadAuthenticator:  notificationsTemplateReadAuthenticator,
		NotificationTemplatesWriteAuthenticator: notificationsTemplateWriteAuthenticator,
		NotificationsManageAuthenticator:        notificationsManageAuthenticator,
	}))
	v1.AddMux(NewNotificationsRouter(NotificationsRouterConfig{
		ErrorWriter:         errorWriter,
		Registrar:           registrar,
		NotificationsFinder: notificationsFinder,

		RequestLogging:                   logging,
		DatabaseAllocator:                databaseAllocator,
		NotificationsWriteAuthenticator:  notificationsWriteAuthenticator,
		NotificationsManageAuthenticator: notificationsManageAuthenticator,
	}))
	v1.AddMux(NewNotifyRouter(NotifyRouterConfig{
		ErrorWriter:          errorWriter,
		Notify:               notify,
		UserStrategy:         userStrategy,
		SpaceStrategy:        spaceStrategy,
		OrganizationStrategy: organizationStrategy,
		EveryoneStrategy:     everyoneStrategy,
		UAAScopeStrategy:     uaaScopeStrategy,
		EmailStrategy:        emailStrategy,

		RequestLogging:                  logging,
		NotificationsWriteAuthenticator: notificationsWriteAuthenticator,
		DatabaseAllocator:               databaseAllocator,
		EmailsWriteAuthenticator:        emailsWriteAuthenticator,
	}))

	// V2
	sendersCollection := collections.NewSendersCollection(models.NewSendersRepository(uuid.NewV4))

	v2 := NewRouterPool()
	v2.AddMux(NewInfoRouter(InfoRouterConfig{
		Version:        2,
		RequestLogging: logging,
	}))
	v2.AddMux(senders.NewRouter(senders.RouterConfig{
		RequestLogging:    logging,
		Authenticator:     notificationsWriteAuthenticator,
		DatabaseAllocator: databaseAllocator,
		SendersCollection: sendersCollection,
	}))
	v2.AddMux(notificationtypes.NewRouter(notificationtypes.RouterConfig{
		NotificationTypesCollection: collections.NotificationTypesCollection{},
	}))

	return VersionRouter{
		1: v1,
		2: v2,
	}
}
