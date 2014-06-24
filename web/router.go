package web

import (
    "log"
    "net/http"
    "strings"

    "github.com/gorilla/mux"
    "github.com/pivotal-cf/cf-notifications/web/handlers"
    "github.com/pivotal-cf/cf-notifications/web/middleware"
)

type Router struct {
    stacks map[string]Stack
}

func NewRouter() Router {
    return Router{
        stacks: map[string]Stack{
            "GET /info": buildUnauthenticatedStack(handlers.NewGetInfo()),
        },
    }
}

func (router Router) Routes() *mux.Router {
    r := mux.NewRouter()
    for methodPath, stack := range router.stacks {
        var name = methodPath
        parts := strings.SplitN(methodPath, " ", 2)
        r.Handle(parts[1], stack).Methods(parts[0]).Name(name)
    }
    return r
}

func buildUnauthenticatedStack(handler http.Handler) Stack {
    logging := middleware.NewLogging(&log.Logger{})
    return NewStack(handler).Use(logging)
}
