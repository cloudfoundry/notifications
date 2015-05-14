package web

import (
	"net/http"

	"github.com/pivotal-golang/lager"
)

type Server struct {
}

func NewServer() Server {
	return Server{}
}

func (s Server) Run(port string, mother MotherInterface, logger lager.Logger) {
	router := NewRouter(mother)
	logger.Info("listen-and-serve", lager.Data{
		"port": port,
	})

	http.ListenAndServe(":"+port, router.Routes())
}
