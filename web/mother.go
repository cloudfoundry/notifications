package web

import (
    "log"
    "os"

    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/cloudfoundry-incubator/notifications/web/handlers"
    "github.com/cloudfoundry-incubator/notifications/web/middleware"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
    "github.com/ryanmoran/stack"
)

type Mother struct {
    logger *log.Logger
    queue  *postal.DeliveryQueue
}

func NewMother() *Mother {
    return &Mother{}
}

func (mother *Mother) Logger() *log.Logger {
    if mother.logger == nil {
        mother.logger = log.New(os.Stdout, "[WEB] ", log.LstdFlags)
    }
    return mother.logger
}

func (mother *Mother) Queue() *postal.DeliveryQueue {
    if mother.queue == nil {
        mother.queue = postal.NewDeliveryQueue()
    }
    return mother.queue
}

func (mother Mother) Courier() postal.Courier {
    env := config.NewEnvironment()
    uaaClient := uaa.NewUAA("", env.UAAHost, env.UAAClientID, env.UAAClientSecret, "")
    uaaClient.VerifySSL = env.VerifySSL
    cloudController := cf.NewCloudController(env.CCHost)

    tokenLoader := postal.NewTokenLoader(&uaaClient)
    userLoader := postal.NewUserLoader(&uaaClient, mother.Logger(), cloudController)
    spaceLoader := postal.NewSpaceLoader(cloudController)
    templateLoader := postal.NewTemplateLoader(postal.NewFileSystem())
    mailer := postal.NewMailer(mother.Queue())

    return postal.NewCourier(tokenLoader, userLoader, spaceLoader, templateLoader, mailer)
}

func (mother Mother) Finder() handlers.Finder {
    clientsRepo, kindsRepo := mother.Repos()
    return handlers.NewFinder(clientsRepo, kindsRepo)
}

func (mother Mother) Repos() (models.ClientsRepo, models.KindsRepo) {
    return models.NewClientsRepo(), models.NewKindsRepo()
}

func (mother Mother) Logging() stack.Middleware {
    return stack.NewLogging(mother.Logger())
}

func (mother Mother) ErrorWriter() handlers.ErrorWriter {
    return handlers.NewErrorWriter()
}

func (mother Mother) Authenticator() middleware.Authenticator {
    return middleware.NewAuthenticator()
}

func (mother Mother) Registrar() handlers.Registrar {
    clientsRepo, kindsRepo := mother.Repos()
    return handlers.NewRegistrar(clientsRepo, kindsRepo)
}
