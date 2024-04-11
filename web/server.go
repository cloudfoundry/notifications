package web

import (
	"database/sql"
	"net/http"

	"fmt"

	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/uaa"
	"github.com/pivotal-golang/lager"
)

type Config struct {
	DBLoggingEnabled     bool
	SkipVerifySSL        bool
	Port                 int
	CORSOrigin           string
	QueueWaitMaxDuration int
	MaxQueueLength       int
	SQLDB                *sql.DB
	Queue                gobble.QueueInterface
	Logger               lager.Logger

	UAATokenValidator *uaa.TokenValidator
	UAAHost           string
	UAAClientID       string
	UAAClientSecret   string
	DefaultUAAScopes  []string
	CCHost            string
}

type Server struct{}

func NewServer() Server {
	return Server{}
}

func (s Server) Run(config Config) {
	config.Logger.Info("listen-and-serve", lager.Data{
		"port": config.Port,
	})

	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), NewRouter(config))
}
