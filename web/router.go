package web

import (
    "net"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/mail"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/cloudfoundry-incubator/notifications/web/handlers"
    "github.com/gorilla/mux"
    "github.com/nu7hatch/gouuid"

    "github.com/ryanmoran/stack"
)

const WorkerCount = 10

type Router struct {
    stacks map[string]stack.Stack
}

func NewRouter() Router {
    mother := NewMother()

    StartWorkers(mother)

    notify := handlers.NewNotify(mother.Courier(), mother.Finder())
    logging := mother.Logging()
    errorWriter := mother.ErrorWriter()
    authenticator := mother.Authenticator()
    registrar := mother.Registrar()

    return Router{
        stacks: map[string]stack.Stack{
            "GET /info":           stack.NewStack(handlers.NewGetInfo()).Use(logging),
            "POST /users/{guid}":  stack.NewStack(handlers.NewNotifyUser(notify, errorWriter)).Use(logging, authenticator),
            "POST /spaces/{guid}": stack.NewStack(handlers.NewNotifySpace(notify, errorWriter)).Use(logging, authenticator),
            "PUT /registration":   stack.NewStack(handlers.NewRegistration(registrar, errorWriter)).Use(logging, authenticator),
        },
    }
}

func StartWorkers(mother *Mother) {
    env := config.NewEnvironment()
    for i := 0; i < WorkerCount; i++ {
        mailClient, err := mail.NewClient(env.SMTPUser, env.SMTPPass, net.JoinHostPort(env.SMTPHost, env.SMTPPort))
        if err != nil {
            panic(err)
        }
        mailClient.Insecure = !env.VerifySSL
        worker := postal.NewDeliveryWorker(uuid.NewV4, mother.Logger(), mailClient, mother.Queue())
        worker.Run()
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
