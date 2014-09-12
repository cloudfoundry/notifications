package mother

import (
    "log"
    "os"

    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/gobble"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/cloudfoundry-incubator/notifications/web/handlers"
    "github.com/cloudfoundry-incubator/notifications/web/middleware"
    "github.com/cloudfoundry-incubator/notifications/web/services"
    "github.com/nu7hatch/gouuid"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
    "github.com/ryanmoran/stack"
)

type Mother struct {
    logger *log.Logger
    queue  *gobble.Queue
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

func (mother *Mother) Queue() *gobble.Queue {
    if mother.queue == nil {
        mother.queue = gobble.NewQueue()
    }
    return mother.queue
}

func (mother Mother) NewUAARecipe() postal.UAARecipe {
    env := config.NewEnvironment()
    uaaClient := uaa.NewUAA("", env.UAAHost, env.UAAClientID, env.UAAClientSecret, "")
    uaaClient.VerifySSL = env.VerifySSL
    cloudController := cf.NewCloudController(env.CCHost)

    tokenLoader := postal.NewTokenLoader(&uaaClient)
    userLoader := postal.NewUserLoader(&uaaClient, mother.Logger(), cloudController)
    spaceLoader := postal.NewSpaceLoader(cloudController)
    templateLoader := postal.NewTemplateLoader(postal.NewFileSystem())
    mailer := mother.Mailer()
    receiptsRepo := models.NewReceiptsRepo()

    return postal.NewUAARecipe(tokenLoader, userLoader, spaceLoader, templateLoader, mailer, receiptsRepo)
}

func (mother Mother) NotificationFinder() services.NotificationFinder {
    clientsRepo, kindsRepo := mother.Repos()
    return services.NewNotificationFinder(clientsRepo, kindsRepo)
}

func (mother Mother) Mailer() postal.Mailer {
    return postal.NewMailer(mother.Queue(), uuid.NewV4, mother.UnsubscribesRepo())
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

func (mother Mother) Authenticator(scopes []string) middleware.Authenticator {
    return middleware.NewAuthenticator(scopes)
}

func (mother Mother) Registrar() services.Registrar {
    clientsRepo, kindsRepo := mother.Repos()
    return services.NewRegistrar(clientsRepo, kindsRepo)
}

func (mother Mother) PreferencesFinder() *services.PreferencesFinder {
    return services.NewPreferencesFinder(models.NewPreferencesRepo())
}

func (mother Mother) PreferenceUpdater() services.PreferenceUpdater {
    return services.NewPreferenceUpdater(mother.UnsubscribesRepo())
}

func (mother Mother) UnsubscribesRepo() models.UnsubscribesRepo {
    return models.NewUnsubscribesRepo()
}
