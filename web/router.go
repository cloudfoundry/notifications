package web

import (
    "log"
    "os"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/web/handlers"
    "github.com/gorilla/mux"
    "github.com/ryanmoran/stack"
)

type Router struct {
    stacks map[string]stack.Stack
}

func NewRouter() Router {
    logger := log.New(os.Stdout, "[WEB] ", log.LstdFlags)
    logging := stack.NewLogging(logger)

    return Router{
        stacks: map[string]stack.Stack{
            "GET /info":          stack.NewStack(handlers.NewGetInfo()).Use(logging),
            "POST /users/{uuid}": stack.NewStack(handlers.NewNotifyUser(logger)).Use(logging),
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
