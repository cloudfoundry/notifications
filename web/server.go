package web

import (
	"net/http"
	"strconv"

	"github.com/pivotal-golang/lager"
)

type Config struct {
	Port             int
	Logger           lager.Logger
	DBLoggingEnabled bool
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
