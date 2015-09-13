package web

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	v1web "github.com/cloudfoundry-incubator/notifications/v1/web"
	v2web "github.com/cloudfoundry-incubator/notifications/v2/web"
)

type MotherInterface interface {
	Queue() gobble.QueueInterface

	EmailStrategy() services.EmailStrategy
	UserStrategy() services.UserStrategy
	SpaceStrategy() services.SpaceStrategy
	OrganizationStrategy() services.OrganizationStrategy
	EveryoneStrategy() services.EveryoneStrategy
	UAAScopeStrategy() services.UAAScopeStrategy
}

func NewRouter(mother MotherInterface, config Config) http.Handler {
	v1 := v1web.NewRouter(NewMuxer(), mother, v1web.Config{
		DBLoggingEnabled: config.DBLoggingEnabled,
		Logger:           config.Logger,
		UAAPublicKey:     config.UAAPublicKey,
		CORSOrigin:       config.CORSOrigin,
		SQLDB:            config.SQLDB,
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
