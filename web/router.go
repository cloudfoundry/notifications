package web

import (
    "log"
    "net/http"
    "strings"

    "github.com/gorilla/mux"
    "github.com/pivotal-cf/cf-notifications/web/handlers"
    "github.com/ryanmoran/stack"
)

type Router struct {
    stacks map[string]stack.Stack
}

func NewRouter() Router {
    return Router{
        stacks: map[string]stack.Stack{
            "GET /info": buildStack(handlers.NewGetInfo()),
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

func buildStack(handler http.Handler) stack.Stack {
    logging := stack.NewLogging(&log.Logger{})
    return stack.NewStack(handler).Use(logging)
}
