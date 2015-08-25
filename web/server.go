package web

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/pivotal-golang/lager"
)

type Config struct {
	DBLoggingEnabled bool
	SkipVerifySSL    bool
	Port             int
	Logger           lager.Logger
	CORSOrigin       string
	SQLDB            *sql.DB

	UAAPublicKey    string
	UAAHost         string
	UAAClientID     string
	UAAClientSecret string
	CCHost          string
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
