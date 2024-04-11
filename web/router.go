package web

import (
	"net/http"

	v1web "github.com/cloudfoundry-incubator/notifications/v1/web"
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
		MaxQueueLength:    config.MaxQueueLength,
	})

	return VersionRouter{
		1: v1,
	}
}
