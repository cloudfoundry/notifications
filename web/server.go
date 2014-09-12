package web

import (
    "log"
    "net/http"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/mother"
)

type Server struct {
}

func NewServer() Server {
    return Server{}
}

func (s Server) Run(mother *mother.Mother) {
    env := config.NewEnvironment()
    router := NewRouter(mother)
    log.Printf("Listening on localhost:%s\n", env.Port)

    http.ListenAndServe(":"+env.Port, router.Routes())
}
