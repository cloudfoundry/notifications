package web

import (
	"database/sql"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/v1/services"
	v1web "github.com/cloudfoundry-incubator/notifications/v1/web"
	v2web "github.com/cloudfoundry-incubator/notifications/v2/web"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/cloudfoundry-incubator/notifications/web/webutil"
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
	return VersionRouter{
		1: v1web.NewRouter(NewMuxer(), mother, v1web.Config{
			DBLoggingEnabled: config.DBLoggingEnabled,
		}),
		2: v2web.NewRouter(NewMuxer(), mother, v2web.Config{
			DBLoggingEnabled: config.DBLoggingEnabled,
		}),
	}
}
