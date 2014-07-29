package web

import (
    "log"
    "net"
    "os"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/mail"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/cloudfoundry-incubator/notifications/web/handlers"
    "github.com/cloudfoundry-incubator/notifications/web/middleware"
    "github.com/gorilla/mux"
    "github.com/nu7hatch/gouuid"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
    "github.com/ryanmoran/stack"
)

type Router struct {
    stacks map[string]stack.Stack
}

func NewRouter() Router {
    logger := log.New(os.Stdout, "[WEB] ", log.LstdFlags)
    logging := stack.NewLogging(logger)

    authenticator := middleware.NewAuthenticator()

    env := config.NewEnvironment()
    mailClient, err := mail.NewClient(env.SMTPUser, env.SMTPPass, net.JoinHostPort(env.SMTPHost, env.SMTPPort))
    if err != nil {
        panic(err)
    }
    mailClient.Insecure = !env.VerifySSL

    uaaClient := uaa.NewUAA("", env.UAAHost, env.UAAClientID, env.UAAClientSecret, "")
    uaaClient.VerifySSL = env.VerifySSL

    cloudController := cf.NewCloudController(env.CCHost)

    userLoader := postal.NewUserLoader(&uaaClient, logger, cloudController)
    spaceLoader := postal.NewSpaceLoader(cloudController)
    fs := postal.NewFileSystem()
    templateLoader := postal.NewTemplateLoader(&fs)
    mailer := postal.NewMailer(uuid.NewV4, logger, &mailClient)

    courier := postal.NewCourier(&uaaClient, userLoader, spaceLoader, templateLoader, mailer)

    return Router{
        stacks: map[string]stack.Stack{
            "GET /info":           stack.NewStack(handlers.NewGetInfo()).Use(logging),
            "POST /users/{guid}":  stack.NewStack(handlers.NewNotifyUser(courier)).Use(logging, authenticator),
            "POST /spaces/{guid}": stack.NewStack(handlers.NewNotifySpace(courier)).Use(logging, authenticator),
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
