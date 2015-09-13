package web

import (
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
}

type mother interface {
	Repos() (models.ClientsRepo, models.KindsRepo)
	NotificationsFinder() services.NotificationsFinder
	EmailStrategy() services.EmailStrategy
	UserStrategy() services.UserStrategy
	SpaceStrategy() services.SpaceStrategy
	OrganizationStrategy() services.OrganizationStrategy
	EveryoneStrategy() services.EveryoneStrategy
	UAAScopeStrategy() services.UAAScopeStrategy
	PreferencesFinder() *services.PreferencesFinder
	PreferenceUpdater() services.PreferenceUpdater
	TemplateServiceObjects() (services.TemplateCreator, services.TemplateFinder, services.TemplateUpdater, services.TemplateDeleter, services.TemplateLister, services.TemplateAssigner, services.TemplateAssociationLister)
	NotificationsUpdater() services.NotificationsUpdater
	MessageFinder() services.MessageFinder
	SQLDatabase() *sql.DB
}

func NewRouter(mx muxer, mom mother, config Config) http.Handler {
	clientsRepo, kindsRepo := mom.Repos()

	registrar := services.NewRegistrar(clientsRepo, kindsRepo)

	notificationsFinder := mom.NotificationsFinder()
	emailStrategy := mom.EmailStrategy()
	userStrategy := mom.UserStrategy()
	spaceStrategy := mom.SpaceStrategy()
	organizationStrategy := mom.OrganizationStrategy()
	everyoneStrategy := mom.EveryoneStrategy()
	uaaScopeStrategy := mom.UAAScopeStrategy()
	notifyObj := notify.NewNotify(mom.NotificationsFinder(), registrar)
	preferencesFinder := mom.PreferencesFinder()
	preferenceUpdater := mom.PreferenceUpdater()
	templateCreator, templateFinder, templateUpdater, templateDeleter, templateLister, templateAssigner, templateAssociationLister := mom.TemplateServiceObjects()
	notificationsUpdater := mom.NotificationsUpdater()
	messageFinder := mom.MessageFinder()
	errorWriter := webutil.NewErrorWriter()

	requestCounter := middleware.NewRequestCounter(mx.GetRouter(), metrics.DefaultLogger)
	logging := middleware.NewRequestLogging(config.Logger)

	auth := func(scope ...string) middleware.Authenticator {
		return middleware.NewAuthenticator(config.UAAPublicKey, scope...)
	}

	databaseAllocator := middleware.NewDatabaseAllocator(mom.SQLDatabase(), config.DBLoggingEnabled)
	cors := middleware.NewCORS(config.CORSOrigin)

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
