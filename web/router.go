package web

import (
	"net/http"

	v1web "github.com/cloudfoundry-incubator/notifications/v1/web"
	v2web "github.com/cloudfoundry-incubator/notifications/v2/web"
)

func NewRouter(config Config) http.Handler {
	v1 := v1web.NewRouter(NewMuxer(), v1web.Config{
		UAATokenValidator: config.UAATokenValidator,
		UAAClientID:       config.UAAClientID,
		UAAClientSecret:   config.UAAClientSecret,
		DefaultUAAScopes:  config.DefaultUAAScopes,
		DBLoggingEnabled:  config.DBLoggingEnabled,
		Logger:            config.Logger,
		VerifySSL:         !config.SkipVerifySSL,
		CCHost:            config.CCHost,
		CORSOrigin:        config.CORSOrigin,
		SQLDB:             config.SQLDB,
	})

	v2 := v2web.NewRouter(NewMuxer(), v2web.Config{
		DBLoggingEnabled:  config.DBLoggingEnabled,
		SkipVerifySSL:     config.SkipVerifySSL,
		SQLDB:             config.SQLDB,
		Logger:            config.Logger,
		Queue:             config.Queue,
		UAATokenValidator: config.UAATokenValidator,
		UAAHost:           config.UAAHost,
		UAAClientID:       config.UAAClientID,
		UAAClientSecret:   config.UAAClientSecret,
		CCHost:            config.CCHost,
	})

	return VersionRouter{
		1: v1,
		2: v2,
	}
}
