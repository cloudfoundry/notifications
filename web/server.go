package web

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/cloudfoundry-incubator/notifications/uaa"
	"github.com/pivotal-golang/lager"
)

type Config struct {
	DBLoggingEnabled     bool
	SkipVerifySSL        bool
	Port                 int
	CORSOrigin           string
	QueueWaitMaxDuration int
	SQLDB                *sql.DB
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

func (s Server) Run(mother MotherInterface, config Config) {
	config.Logger.Info("listen-and-serve", lager.Data{
		"port": config.Port,
	})

	http.ListenAndServe(":"+strconv.Itoa(config.Port), NewRouter(mother, config))
}
