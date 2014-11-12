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

func (s Server) Run(port string, mother MotherInterface) {
	router := NewRouter(mother)
	log.Printf("Listening on localhost:%s\n", port)

	http.ListenAndServe(":"+port, router.Routes())
}
