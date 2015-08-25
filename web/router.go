package web

import (
	"database/sql"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	v1web "github.com/cloudfoundry-incubator/notifications/v1/web"
	v2web "github.com/cloudfoundry-incubator/notifications/v2/web"
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
	SQLDatabase() *sql.DB
	Queue() gobble.QueueInterface
}

func NewRouter(mother MotherInterface, config Config) http.Handler {
	v1 := v1web.NewRouter(NewMuxer(), mother, v1web.Config{
		DBLoggingEnabled: config.DBLoggingEnabled,
		Logger:           config.Logger,
		UAAPublicKey:     config.UAAPublicKey,
		CORSOrigin:       config.CORSOrigin,
	})

	v2 := v2web.NewRouter(NewMuxer(), v2web.Config{
		DBLoggingEnabled: config.DBLoggingEnabled,
		SkipVerifySSL:    config.SkipVerifySSL,
		SQLDB:            config.SQLDB,
		Logger:           config.Logger,
		Queue:            mother.Queue(),
		UAAPublicKey:     config.UAAPublicKey,
		UAAHost:          config.UAAHost,
		UAAClientID:      config.UAAClientID,
		UAAClientSecret:  config.UAAClientSecret,
		CCHost:           config.CCHost,
	})

	return VersionRouter{
		1: v1,
		2: v2,
	}
}
