package web

import (
	"database/sql"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/cloudfoundry-incubator/notifications/web/v1/clients"
	"github.com/cloudfoundry-incubator/notifications/web/v1/info"
	"github.com/cloudfoundry-incubator/notifications/web/v1/messages"
	"github.com/cloudfoundry-incubator/notifications/web/v1/notifications"
	"github.com/cloudfoundry-incubator/notifications/web/v1/notify"
	"github.com/cloudfoundry-incubator/notifications/web/v1/preferences"
	"github.com/cloudfoundry-incubator/notifications/web/v1/templates"
	"github.com/cloudfoundry-incubator/notifications/web/v2/notificationtypes"
	"github.com/cloudfoundry-incubator/notifications/web/v2/senders"
	"github.com/cloudfoundry-incubator/notifications/web/webutil"
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
	Logging() middleware.RequestLogging
	ErrorWriter() webutil.ErrorWriter
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
	notifyObj := notify.NewNotify(mother.NotificationsFinder(), registrar)
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

	v1 := NewRouterPool()
	v1.AddMux(info.NewRouter(info.RouterConfig{
		Version:        1,
		RequestLogging: logging,
	}))
	v1.AddMux(preferences.NewRouter(preferences.RouterConfig{
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
	v1.AddMux(clients.NewRouter(clients.RouterConfig{
		ErrorWriter:      errorWriter,
		TemplateAssigner: templateAssigner,

		RequestLogging:                   logging,
		DatabaseAllocator:                databaseAllocator,
		NotificationsManageAuthenticator: notificationsManageAuthenticator,
	}))
	v1.AddMux(messages.NewRouter(messages.RouterConfig{
		ErrorWriter:   errorWriter,
		MessageFinder: messageFinder,

		RequestLogging:                               logging,
		DatabaseAllocator:                            databaseAllocator,
		NotificationsWriteOrEmailsWriteAuthenticator: notificationsWriteOrEmailsWriteAuthenticator,
	}))
	v1.AddMux(templates.NewRouter(templates.RouterConfig{
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
	v1.AddMux(notifications.NewRouter(notifications.RouterConfig{
		ErrorWriter:          errorWriter,
		Registrar:            registrar,
		NotificationsFinder:  notificationsFinder,
		NotificationsUpdater: notificationsUpdater,
		TemplateAssigner:     templateAssigner,

		RequestLogging:                   logging,
		DatabaseAllocator:                databaseAllocator,
		NotificationsWriteAuthenticator:  notificationsWriteAuthenticator,
		NotificationsManageAuthenticator: notificationsManageAuthenticator,
	}))
	v1.AddMux(notify.NewRouter(notify.RouterConfig{
		ErrorWriter:          errorWriter,
		Notify:               notifyObj,
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
	sendersRepository := models.NewSendersRepository(uuid.NewV4)
	sendersCollection := collections.NewSendersCollection(sendersRepository)
	notificationTypesCollection := collections.NewNotificationTypesCollection(models.NewNotificationTypesRepository(uuid.NewV4), sendersRepository)

	v2 := NewRouterPool()
	v2.AddMux(info.NewRouter(info.RouterConfig{
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
		RequestLogging:              logging,
		Authenticator:               notificationsWriteAuthenticator,
		DatabaseAllocator:           databaseAllocator,
		NotificationTypesCollection: notificationTypesCollection,
	}))

	return VersionRouter{
		1: v1,
		2: v2,
	}
}
