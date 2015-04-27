package web

import (
	"log"
	"net/http"
)

type Server struct {
}

func NewServer() Server {
	return Server{}
}

func (s Server) Run(port string, mother MotherInterface, logger *log.Logger) {
	router := NewRouter(mother)
	logger.Printf("Listening on localhost:%s\n", port)

	http.ListenAndServe(":"+port, router.Routes())
}
