package web

import (
	"crypto/rand"
	"database/sql"
	"net/http"
	"time"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/uaa"
	"github.com/cloudfoundry-incubator/notifications/util"
	"github.com/cloudfoundry-incubator/notifications/v1/collections"
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
	"github.com/gorilla/mux"
	"github.com/pivotal-golang/lager"
	"github.com/rcrowley/go-metrics"
	"github.com/rcrowley/go-metrics/exp"
	"github.com/ryanmoran/stack"
)

type muxer interface {
	Handle(method, path string, handler stack.Handler, middleware ...stack.Middleware)
	GetRouter() *mux.Router
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

type Config struct {
	UAATokenValidator    *uaa.TokenValidator
	UAAClientID          string
	UAAClientSecret      string
	DefaultUAAScopes     []string
	VerifySSL            bool
	CCHost               string
	DBLoggingEnabled     bool
	Logger               lager.Logger
	CORSOrigin           string
	SQLDB                *sql.DB
	QueueWaitMaxDuration int
	MaxQueueLength       int
}

func NewRouter(mx muxer, config Config) http.Handler {
	guidGenerator := util.NewIDGenerator(rand.Reader)
	clock := util.NewClock()

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

	templatesCollection := collections.NewTemplatesCollection(clientsRepo, kindsRepo, templatesRepo)

	templateFinder := services.NewTemplateFinder(templatesRepo)
	templateUpdater := services.NewTemplateUpdater(templatesRepo)
	templateLister := services.NewTemplateLister(templatesRepo)

	notifyObj := notify.NewNotify(notificationsFinder, registrar)

	gobbleQueue := gobble.NewQueue(gobble.NewDatabase(config.SQLDB), clock, gobble.Config{
		WaitMaxDuration: time.Duration(config.QueueWaitMaxDuration) * time.Millisecond,
		MaxQueueLength:  config.MaxQueueLength,
	})

	v1enqueuer := services.NewEnqueuer(gobbleQueue, messagesRepo, gobble.Initializer{})

	uaaClient := uaa.NewZonedUAAClient(config.UAAClientID, config.UAAClientSecret, config.VerifySSL, config.UAATokenValidator)
	cloudController := cf.NewCloudController(config.CCHost, !config.VerifySSL)
	tokenLoader := uaa.NewTokenLoader(uaaClient)
	spaceLoader := services.NewSpaceLoader(cloudController)
	organizationLoader := services.NewOrganizationLoader(cloudController)
	findsUserIDs := services.NewFindsUserIDs(cloudController, uaaClient)
	allUsers := services.NewAllUsers(uaaClient)

	emailStrategy := services.NewEmailStrategy(v1enqueuer)
	userStrategy := services.NewUserStrategy(v1enqueuer)
	spaceStrategy := services.NewSpaceStrategy(tokenLoader, spaceLoader, organizationLoader, findsUserIDs, v1enqueuer)
	organizationStrategy := services.NewOrganizationStrategy(tokenLoader, organizationLoader, findsUserIDs, v1enqueuer)
	everyoneStrategy := services.NewEveryoneStrategy(tokenLoader, allUsers, v1enqueuer)
	uaaScopeStrategy := services.NewUAAScopeStrategy(tokenLoader, findsUserIDs, v1enqueuer, config.DefaultUAAScopes)

	errorWriter := webutil.NewErrorWriter()

	requestCounter := middleware.NewRequestCounter(mx.GetRouter())
	requestLogging := middleware.NewRequestLogging(config.Logger, clock)
	databaseAllocator := middleware.NewDatabaseAllocator(config.SQLDB, config.DBLoggingEnabled)
	cors := middleware.NewCORS(config.CORSOrigin)
	auth := func(scope ...string) middleware.Authenticator {
		return middleware.NewAuthenticator(config.UAATokenValidator, scope...)
	}

	mx.GetRouter().Handle("/debug/metrics", exp.ExpHandler(metrics.DefaultRegistry)).Methods("GET")

	info.Routes{
		RequestCounter: requestCounter,
		RequestLogging: requestLogging,
	}.Register(mx)

	preferences.Routes{
		CORS:                                     cors,
		RequestCounter:                           requestCounter,
		RequestLogging:                           requestLogging,
		DatabaseAllocator:                        databaseAllocator,
		NotificationPreferencesReadAuthenticator: auth("notification_preferences.read"),
		NotificationPreferencesWriteAuthenticator: auth("notification_preferences.write"),
		NotificationPreferencesAdminAuthenticator: auth("notification_preferences.admin"),

		ErrorWriter:       errorWriter,
		PreferencesFinder: preferencesFinder,
		PreferenceUpdater: preferenceUpdater,
	}.Register(mx)

	clients.Routes{
		RequestCounter:                   requestCounter,
		RequestLogging:                   requestLogging,
		DatabaseAllocator:                databaseAllocator,
		NotificationsManageAuthenticator: auth("notifications.manage"),

		ErrorWriter:      errorWriter,
		TemplateAssigner: templatesCollection,
	}.Register(mx)

	messages.Routes{
		RequestCounter:    requestCounter,
		RequestLogging:    requestLogging,
		DatabaseAllocator: databaseAllocator,
		NotificationsWriteOrEmailsWriteAuthenticator: auth("notifications.write", "emails.write"),

		ErrorWriter:   errorWriter,
		MessageFinder: messageFinder,
	}.Register(mx)

	templates.Routes{
		RequestCounter:                          requestCounter,
		RequestLogging:                          requestLogging,
		DatabaseAllocator:                       databaseAllocator,
		NotificationTemplatesReadAuthenticator:  auth("notification_templates.read"),
		NotificationTemplatesWriteAuthenticator: auth("notification_templates.write"),
		NotificationsManageAuthenticator:        auth("notifications.manage"),

		ErrorWriter:               errorWriter,
		TemplateFinder:            templateFinder,
		TemplateUpdater:           templateUpdater,
		TemplateCreator:           templatesCollection,
		TemplateDeleter:           templatesCollection,
		TemplateLister:            templateLister,
		TemplateAssociationLister: templatesCollection,
	}.Register(mx)

	notifications.Routes{
		RequestCounter:                   requestCounter,
		RequestLogging:                   requestLogging,
		DatabaseAllocator:                databaseAllocator,
		NotificationsWriteAuthenticator:  auth("notifications.write"),
		NotificationsManageAuthenticator: auth("notifications.manage"),

		ErrorWriter:          errorWriter,
		Registrar:            registrar,
		NotificationsFinder:  notificationsFinder,
		NotificationsUpdater: notificationsUpdater,
		TemplateAssigner:     templatesCollection,
	}.Register(mx)

	notify.Routes{
		RequestCounter:                  requestCounter,
		RequestLogging:                  requestLogging,
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
