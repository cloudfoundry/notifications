package web

import (
    "log"
    "net/http"
    "os"
)

type Server struct {
}

func NewServer() Server {
    return Server{}
}

func (s Server) Run() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "3000"
    }
    router := NewRouter()
    log.Printf("Listening on localhost:%s\n", port)

    http.ListenAndServe(":"+port, router.Routes())
}
