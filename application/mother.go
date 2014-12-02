package application

import (
	"log"
	"os"
	"path"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/config"
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/mail"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/postal/strategies"
	"github.com/cloudfoundry-incubator/notifications/postal/utilities"
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

func (mother Mother) UserStrategy() strategies.UserStrategy {
	env := config.NewEnvironment()
	finder := mother.TemplateFinder()
	uaaClient := uaa.NewUAA("", env.UAAHost, env.UAAClientID, env.UAAClientSecret, "")
	uaaClient.VerifySSL = env.VerifySSL

	tokenLoader := utilities.NewTokenLoader(&uaaClient)
	userLoader := utilities.NewUserLoader(&uaaClient, mother.Logger())
	templatesLoader := utilities.NewTemplatesLoader(finder)
	mailer := mother.Mailer()
	receiptsRepo := models.NewReceiptsRepo()

	return strategies.NewUserStrategy(tokenLoader, userLoader, templatesLoader, mailer, receiptsRepo)
}

func (mother Mother) SpaceStrategy() strategies.SpaceStrategy {
	env := config.NewEnvironment()
	finder := mother.TemplateFinder()
	uaaClient := uaa.NewUAA("", env.UAAHost, env.UAAClientID, env.UAAClientSecret, "")
	uaaClient.VerifySSL = env.VerifySSL
	cloudController := cf.NewCloudController(env.CCHost, !env.VerifySSL)

	tokenLoader := utilities.NewTokenLoader(&uaaClient)
	userLoader := utilities.NewUserLoader(&uaaClient, mother.Logger())
	spaceLoader := utilities.NewSpaceLoader(cloudController)
	organizationLoader := utilities.NewOrganizationLoader(cloudController)
	templatesLoader := utilities.NewTemplatesLoader(finder)
	mailer := mother.Mailer()
	receiptsRepo := models.NewReceiptsRepo()
	findsUserGUIDs := utilities.NewFindsUserGUIDs(cloudController, &uaaClient)

	return strategies.NewSpaceStrategy(tokenLoader, userLoader, spaceLoader, organizationLoader, findsUserGUIDs, templatesLoader, mailer, receiptsRepo)
}

func (mother Mother) OrganizationStrategy() strategies.OrganizationStrategy {
	env := config.NewEnvironment()
	finder := mother.TemplateFinder()
	uaaClient := uaa.NewUAA("", env.UAAHost, env.UAAClientID, env.UAAClientSecret, "")
	uaaClient.VerifySSL = env.VerifySSL
	cloudController := cf.NewCloudController(env.CCHost, !env.VerifySSL)

	tokenLoader := utilities.NewTokenLoader(&uaaClient)
	userLoader := utilities.NewUserLoader(&uaaClient, mother.Logger())
	organizationLoader := utilities.NewOrganizationLoader(cloudController)
	findsUserGUIDs := utilities.NewFindsUserGUIDs(cloudController, &uaaClient)
	templatesLoader := utilities.NewTemplatesLoader(finder)
	mailer := mother.Mailer()
	receiptsRepo := models.NewReceiptsRepo()

	return strategies.NewOrganizationStrategy(tokenLoader, userLoader, organizationLoader, findsUserGUIDs, templatesLoader, mailer, receiptsRepo)
}

func (mother Mother) EveryoneStrategy() strategies.EveryoneStrategy {
	env := config.NewEnvironment()
	finder := mother.TemplateFinder()
	uaaClient := uaa.NewUAA("", env.UAAHost, env.UAAClientID, env.UAAClientSecret, "")
	uaaClient.VerifySSL = env.VerifySSL
	tokenLoader := utilities.NewTokenLoader(&uaaClient)
	allUsers := utilities.NewAllUsers(&uaaClient)

	templatesLoader := utilities.NewTemplatesLoader(finder)
	mailer := mother.Mailer()
	receiptsRepo := models.NewReceiptsRepo()

	return strategies.NewEveryoneStrategy(tokenLoader, allUsers, templatesLoader, mailer, receiptsRepo)
}

func (mother Mother) UAAScopeStrategy() strategies.UAAScopeStrategy {
	env := config.NewEnvironment()
	finder := mother.TemplateFinder()
	uaaClient := uaa.NewUAA("", env.UAAHost, env.UAAClientID, env.UAAClientSecret, "")
	uaaClient.VerifySSL = env.VerifySSL
	cloudController := cf.NewCloudController(env.CCHost, !env.VerifySSL)

	tokenLoader := utilities.NewTokenLoader(&uaaClient)
	userLoader := utilities.NewUserLoader(&uaaClient, mother.Logger())
	findsUserGUIDs := utilities.NewFindsUserGUIDs(cloudController, &uaaClient)
	templatesLoader := utilities.NewTemplatesLoader(finder)
	receiptsRepo := models.NewReceiptsRepo()
	mailer := mother.Mailer()

	return strategies.NewUAAScopeStrategy(tokenLoader, userLoader, findsUserGUIDs, templatesLoader, mailer, receiptsRepo)
}

func (mother Mother) FileSystem() services.FileSystemInterface {
	return postal.NewFileSystem()
}

func (mother Mother) EmailStrategy() strategies.EmailStrategy {
	return strategies.NewEmailStrategy(mother.Mailer(), utilities.NewTemplatesLoader(mother.TemplateFinder()))
}

func (mother Mother) NotificationsFinder() services.NotificationsFinder {
	clientsRepo, kindsRepo := mother.Repos()
	return services.NewNotificationsFinder(clientsRepo, kindsRepo, mother.Database())
}

func (mother Mother) Mailer() strategies.Mailer {
	return strategies.NewMailer(mother.Queue(), uuid.NewV4)
}

func (mother Mother) MailClient() *mail.Client {
	env := config.NewEnvironment()
	mailConfig := mail.Config{
		User:           env.SMTPUser,
		Pass:           env.SMTPPass,
		Host:           env.SMTPHost,
		Port:           env.SMTPPort,
		TestMode:       env.TestMode,
		SkipVerifySSL:  env.VerifySSL,
		DisableTLS:     !env.SMTPTLS,
		LoggingEnabled: env.SMTPLoggingEnabled,
	}
	client, err := mail.NewClient(mailConfig, mother.Logger())
	if err != nil {
		panic(err)
	}

	return client
}

func (mother Mother) Repos() (models.ClientsRepo, models.KindsRepo) {
	return models.NewClientsRepo(), mother.KindsRepo()
}

func (mother Mother) Logging() stack.Middleware {
	return stack.NewLogging(mother.Logger())
}

func (mother Mother) ErrorWriter() handlers.ErrorWriter {
	return handlers.NewErrorWriter()
}

func (mother Mother) Authenticator(scopes ...string) middleware.Authenticator {
	return middleware.NewAuthenticator(config.UAAPublicKey, scopes...)
}

func (mother Mother) Registrar() services.Registrar {
	clientsRepo, kindsRepo := mother.Repos()
	return services.NewRegistrar(clientsRepo, kindsRepo)
}

func (mother Mother) Database() models.DatabaseInterface {
	env := config.NewEnvironment()
	migrationsPath := path.Join(env.RootPath, env.ModelMigrationsDir)
	return models.NewDatabase(env.DatabaseURL, migrationsPath)
}

func (mother Mother) PreferencesFinder() *services.PreferencesFinder {
	return services.NewPreferencesFinder(models.NewPreferencesRepo(), mother.GlobalUnsubscribesRepo(), mother.Database())
}

func (mother Mother) PreferenceUpdater() services.PreferenceUpdater {
	return services.NewPreferenceUpdater(mother.GlobalUnsubscribesRepo(), mother.UnsubscribesRepo(), mother.KindsRepo())
}

func (mother Mother) TemplateFinder() services.TemplateFinder {
	env := config.NewEnvironment()
	database := mother.Database()
	repo := mother.TemplatesRepo()
	fileSystem := mother.FileSystem()

	return services.NewTemplateFinder(repo, env.RootPath, database, fileSystem)
}

func (mother Mother) TemplateServiceObjects() (services.TemplateFinder, services.TemplateUpdater, services.TemplateDeleter) {
	env := config.NewEnvironment()
	database := mother.Database()
	repo := mother.TemplatesRepo()
	fileSystem := mother.FileSystem()

	return services.NewTemplateFinder(repo, env.RootPath, database, fileSystem),
		services.NewTemplateUpdater(repo, database),
		services.NewTemplateDeleter(repo, database)
}

func (mother Mother) KindsRepo() models.KindsRepo {
	return models.NewKindsRepo()
}

func (mother Mother) UnsubscribesRepo() models.UnsubscribesRepo {
	return models.NewUnsubscribesRepo()
}

func (mother Mother) GlobalUnsubscribesRepo() models.GlobalUnsubscribesRepo {
	return models.NewGlobalUnsubscribesRepo()
}

func (mother Mother) TemplatesRepo() models.TemplatesRepo {
	return models.NewTemplatesRepo()
}

func (mother Mother) CORS() middleware.CORS {
	env := config.NewEnvironment()
	return middleware.NewCORS(env.CORSOrigin)
}
