package web

import (
	"crypto/rand"
	"database/sql"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/cloudfoundry-incubator/notifications/v1/web/clients"
	"github.com/cloudfoundry-incubator/notifications/v1/web/info"
	"github.com/cloudfoundry-incubator/notifications/v1/web/messages"
	"github.com/cloudfoundry-incubator/notifications/v1/web/middleware"
	"github.com/cloudfoundry-incubator/notifications/v1/web/notifications"
	"github.com/cloudfoundry-incubator/notifications/v1/web/notify"
	"github.com/cloudfoundry-incubator/notifications/v1/web/preferences"
	"github.com/cloudfoundry-incubator/notifications/v1/web/templates"
	"github.com/cloudfoundry-incubator/notifications/v1/web/webutil"
	v2models "github.com/cloudfoundry-incubator/notifications/v2/models"
	"github.com/gorilla/mux"
	"github.com/pivotal-golang/lager"
	"github.com/ryanmoran/stack"
)

type muxer interface {
	Handle(method, path string, handler stack.Handler, middleware ...stack.Middleware)
	GetRouter() *mux.Router
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

type Config struct {
	DBLoggingEnabled bool
	Logger           lager.Logger
	UAAPublicKey     string
	CORSOrigin       string
	SQLDB            *sql.DB
}

type mother interface {
	EmailStrategy() services.EmailStrategy
	UserStrategy() services.UserStrategy
	SpaceStrategy() services.SpaceStrategy
	OrganizationStrategy() services.OrganizationStrategy
	EveryoneStrategy() services.EveryoneStrategy
	UAAScopeStrategy() services.UAAScopeStrategy
}

func NewRouter(mx muxer, mom mother, config Config) http.Handler {
	guidGenerator := v2models.NewGUIDGenerator(rand.Reader)

	clientsRepo := models.NewClientsRepo()
	kindsRepo := models.NewKindsRepo()
	globalUnsubscribesRepo := models.NewGlobalUnsubscribesRepo()
	preferencesRepo := models.NewPreferencesRepo()
	unsubscribesRepo := models.NewUnsubscribesRepo()
	messagesRepo := models.NewMessagesRepo(guidGenerator.Generate)
	templatesRepo := models.NewTemplatesRepo()

	registrar := services.NewRegistrar(clientsRepo, kindsRepo)
	notificationsFinder := services.NewNotificationsFinder(clientsRepo, kindsRepo)
	preferencesFinder := services.NewPreferencesFinder(preferencesRepo, globalUnsubscribesRepo)
	preferenceUpdater := services.NewPreferenceUpdater(globalUnsubscribesRepo, unsubscribesRepo, kindsRepo)
	notificationsUpdater := services.NewNotificationsUpdater(kindsRepo)
	messageFinder := services.NewMessageFinder(messagesRepo)

	templateCreator := services.NewTemplateCreator(templatesRepo)
	templateFinder := services.NewTemplateFinder(templatesRepo)
	templateUpdater := services.NewTemplateUpdater(templatesRepo)
	templateDeleter := services.NewTemplateDeleter(templatesRepo)
	templateLister := services.NewTemplateLister(templatesRepo)
	templateAssigner := services.NewTemplateAssigner(clientsRepo, kindsRepo, templatesRepo)
	templateAssociationLister := services.NewTemplateAssociationLister(clientsRepo, kindsRepo, templatesRepo)

	notifyObj := notify.NewNotify(notificationsFinder, registrar)

	emailStrategy := mom.EmailStrategy()
	userStrategy := mom.UserStrategy()
	spaceStrategy := mom.SpaceStrategy()
	organizationStrategy := mom.OrganizationStrategy()
	everyoneStrategy := mom.EveryoneStrategy()
	uaaScopeStrategy := mom.UAAScopeStrategy()

	errorWriter := webutil.NewErrorWriter()

	requestCounter := middleware.NewRequestCounter(mx.GetRouter(), metrics.DefaultLogger)
	logging := middleware.NewRequestLogging(config.Logger)
	databaseAllocator := middleware.NewDatabaseAllocator(config.SQLDB, config.DBLoggingEnabled)
	cors := middleware.NewCORS(config.CORSOrigin)
	auth := func(scope ...string) middleware.Authenticator {
		return middleware.NewAuthenticator(config.UAAPublicKey, scope...)
	}

	info.Routes{
		RequestCounter: requestCounter,
		RequestLogging: logging,
	}.Register(mx)

	preferences.Routes{
		CORS:                                      cors,
		RequestCounter:                            requestCounter,
		RequestLogging:                            logging,
		DatabaseAllocator:                         databaseAllocator,
		NotificationPreferencesReadAuthenticator:  auth("notification_preferences.read"),
		NotificationPreferencesWriteAuthenticator: auth("notification_preferences.write"),
		NotificationPreferencesAdminAuthenticator: auth("notification_preferences.admin"),

		ErrorWriter:       errorWriter,
		PreferencesFinder: preferencesFinder,
		PreferenceUpdater: preferenceUpdater,
	}.Register(mx)

	clients.Routes{
		RequestCounter:                   requestCounter,
		RequestLogging:                   logging,
		DatabaseAllocator:                databaseAllocator,
		NotificationsManageAuthenticator: auth("notifications.manage"),

		ErrorWriter:      errorWriter,
		TemplateAssigner: templateAssigner,
	}.Register(mx)

	messages.Routes{
		RequestCounter:                               requestCounter,
		RequestLogging:                               logging,
		DatabaseAllocator:                            databaseAllocator,
		NotificationsWriteOrEmailsWriteAuthenticator: auth("notifications.write", "emails.write"),

		ErrorWriter:   errorWriter,
		MessageFinder: messageFinder,
	}.Register(mx)

	templates.Routes{
		RequestCounter:                          requestCounter,
		RequestLogging:                          logging,
		DatabaseAllocator:                       databaseAllocator,
		NotificationTemplatesReadAuthenticator:  auth("notification_templates.read"),
		NotificationTemplatesWriteAuthenticator: auth("notification_templates.write"),
		NotificationsManageAuthenticator:        auth("notifications.manage"),

		ErrorWriter:               errorWriter,
		TemplateFinder:            templateFinder,
		TemplateUpdater:           templateUpdater,
		TemplateCreator:           templateCreator,
		TemplateDeleter:           templateDeleter,
		TemplateLister:            templateLister,
		TemplateAssociationLister: templateAssociationLister,
	}.Register(mx)

	notifications.Routes{
		RequestCounter:                   requestCounter,
		RequestLogging:                   logging,
		DatabaseAllocator:                databaseAllocator,
		NotificationsWriteAuthenticator:  auth("notifications.write"),
		NotificationsManageAuthenticator: auth("notifications.manage"),

		ErrorWriter:          errorWriter,
		Registrar:            registrar,
		NotificationsFinder:  notificationsFinder,
		NotificationsUpdater: notificationsUpdater,
		TemplateAssigner:     templateAssigner,
	}.Register(mx)

	notify.Routes{
		RequestCounter:                  requestCounter,
		RequestLogging:                  logging,
		DatabaseAllocator:               databaseAllocator,
		NotificationsWriteAuthenticator: auth("notifications.write"),
		EmailsWriteAuthenticator:        auth("emails.write"),

		ErrorWriter:          errorWriter,
		Notify:               notifyObj,
		UserStrategy:         userStrategy,
		SpaceStrategy:        spaceStrategy,
		OrganizationStrategy: organizationStrategy,
		EveryoneStrategy:     everyoneStrategy,
		UAAScopeStrategy:     uaaScopeStrategy,
		EmailStrategy:        emailStrategy,
	}.Register(mx)

	return mx
}
