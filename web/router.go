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
	"github.com/cloudfoundry-incubator/notifications/web/v2/campaigntypes"
	"github.com/cloudfoundry-incubator/notifications/web/v2/senders"
	templatesv2 "github.com/cloudfoundry-incubator/notifications/web/v2/templates"
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

	v1 := NewMuxer()

	info.Routes{
		Version:        1,
		RequestLogging: logging,
	}.Register(v1)

	preferences.Routes{
		ErrorWriter:       errorWriter,
		PreferencesFinder: preferencesFinder,
		PreferenceUpdater: preferenceUpdater,

		CORS:                                      cors,
		RequestLogging:                            logging,
		DatabaseAllocator:                         databaseAllocator,
		NotificationPreferencesReadAuthenticator:  notificationPreferencesReadAuthenticator,
		NotificationPreferencesWriteAuthenticator: notificationPreferencesWriteAuthenticator,
		NotificationPreferencesAdminAuthenticator: notificationPreferencesAdminAuthenticator,
	}.Register(v1)

	clients.Routes{
		ErrorWriter:      errorWriter,
		TemplateAssigner: templateAssigner,

		RequestLogging:                   logging,
		DatabaseAllocator:                databaseAllocator,
		NotificationsManageAuthenticator: notificationsManageAuthenticator,
	}.Register(v1)

	messages.Routes{
		ErrorWriter:   errorWriter,
		MessageFinder: messageFinder,

		RequestLogging:                               logging,
		DatabaseAllocator:                            databaseAllocator,
		NotificationsWriteOrEmailsWriteAuthenticator: notificationsWriteOrEmailsWriteAuthenticator,
	}.Register(v1)

	templates.Routes{
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
	}.Register(v1)

	notifications.Routes{
		ErrorWriter:          errorWriter,
		Registrar:            registrar,
		NotificationsFinder:  notificationsFinder,
		NotificationsUpdater: notificationsUpdater,
		TemplateAssigner:     templateAssigner,

		RequestLogging:                   logging,
		DatabaseAllocator:                databaseAllocator,
		NotificationsWriteAuthenticator:  notificationsWriteAuthenticator,
		NotificationsManageAuthenticator: notificationsManageAuthenticator,
	}.Register(v1)

	notify.Routes{
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
	}.Register(v1)

	// V2
	sendersRepository := models.NewSendersRepository(uuid.NewV4)
	campaignTypesRepository := models.NewCampaignTypesRepository(uuid.NewV4)

	sendersCollection := collections.NewSendersCollection(sendersRepository)
	campaignTypesCollection := collections.NewCampaignTypesCollection(campaignTypesRepository, sendersRepository)

	v2 := NewMuxer()

	info.Routes{
		Version:        2,
		RequestLogging: logging,
	}.Register(v2)

	senders.Routes{
		RequestLogging:    logging,
		Authenticator:     notificationsWriteAuthenticator,
		DatabaseAllocator: databaseAllocator,
		SendersCollection: sendersCollection,
	}.Register(v2)

	campaigntypes.Routes{
		RequestLogging:          logging,
		Authenticator:           notificationsWriteAuthenticator,
		DatabaseAllocator:       databaseAllocator,
		CampaignTypesCollection: campaignTypesCollection,
	}.Register(v2)

	templatesv2.Routes{
		RequestLogging:      logging,
		Authenticator:       notificationsWriteAuthenticator,
		DatabaseAllocator:   databaseAllocator,
		TemplatesCollection: "",
	}.Register(v2)

	return VersionRouter{
		1: v1,
		2: v2,
	}
}
