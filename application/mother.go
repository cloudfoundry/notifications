package application

import (
	"database/sql"
	"log"
	"os"
	"path"
	"sync"
	"time"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/mail"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/nu7hatch/gouuid"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"
	"github.com/pivotal-golang/lager"
)

type Mother struct {
	uaaClient *uaa.UAA
	sqlDB     *sql.DB
	mutex     sync.Mutex
}

func NewMother() *Mother {
	return &Mother{}
}

func (m *Mother) GobbleDatabase() gobble.DatabaseInterface {
	return gobble.NewDatabase(m.SQLDatabase())
}

func (m *Mother) Queue() gobble.QueueInterface {
	env := NewEnvironment()

	return gobble.NewQueue(m.GobbleDatabase(), gobble.Config{
		WaitMaxDuration: time.Duration(env.GobbleWaitMaxDuration) * time.Millisecond,
	})
}

func (m Mother) UserStrategy() services.UserStrategy {
	return services.NewUserStrategy(m.Enqueuer())
}

func (m Mother) UAAClient() *uaa.UAA {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.uaaClient == nil {
		env := NewEnvironment()
		client := uaa.NewUAA("", env.UAAHost, env.UAAClientID, env.UAAClientSecret, "")
		client.VerifySSL = env.VerifySSL
		m.uaaClient = &client
	}

	return m.uaaClient
}

func (m Mother) SpaceStrategy() services.SpaceStrategy {
	env := NewEnvironment()
	uaaClient := m.UAAClient()
	cloudController := cf.NewCloudController(env.CCHost, !env.VerifySSL)

	tokenLoader := postal.NewTokenLoader(uaaClient)
	spaceLoader := services.NewSpaceLoader(cloudController)
	organizationLoader := services.NewOrganizationLoader(cloudController)
	enqueuer := m.Enqueuer()
	findsUserGUIDs := services.NewFindsUserGUIDs(cloudController, uaaClient)

	return services.NewSpaceStrategy(tokenLoader, spaceLoader, organizationLoader, findsUserGUIDs, enqueuer)
}

func (m Mother) OrganizationStrategy() services.OrganizationStrategy {
	env := NewEnvironment()
	uaaClient := m.UAAClient()
	cloudController := cf.NewCloudController(env.CCHost, !env.VerifySSL)

	tokenLoader := postal.NewTokenLoader(uaaClient)
	organizationLoader := services.NewOrganizationLoader(cloudController)
	findsUserGUIDs := services.NewFindsUserGUIDs(cloudController, uaaClient)
	enqueuer := m.Enqueuer()

	return services.NewOrganizationStrategy(tokenLoader, organizationLoader, findsUserGUIDs, enqueuer)
}

func (m Mother) EveryoneStrategy() services.EveryoneStrategy {
	uaaClient := m.UAAClient()
	tokenLoader := postal.NewTokenLoader(uaaClient)
	allUsers := services.NewAllUsers(uaaClient)
	enqueuer := m.Enqueuer()

	return services.NewEveryoneStrategy(tokenLoader, allUsers, enqueuer)
}

func (m Mother) UAAScopeStrategy() services.UAAScopeStrategy {
	env := NewEnvironment()
	uaaClient := m.UAAClient()
	cloudController := cf.NewCloudController(env.CCHost, !env.VerifySSL)

	tokenLoader := postal.NewTokenLoader(uaaClient)
	findsUserGUIDs := services.NewFindsUserGUIDs(cloudController, uaaClient)
	enqueuer := m.Enqueuer()

	return services.NewUAAScopeStrategy(tokenLoader, findsUserGUIDs, enqueuer)
}

func (m Mother) EmailStrategy() services.EmailStrategy {
	return services.NewEmailStrategy(m.Enqueuer())
}

func (m Mother) NotificationsFinder() services.NotificationsFinder {
	clientsRepo, kindsRepo := m.Repos()
	return services.NewNotificationsFinder(clientsRepo, kindsRepo)
}
func (m Mother) NotificationsUpdater() services.NotificationsUpdater {
	_, kindsRepo := m.Repos()
	return services.NewNotificationsUpdater(kindsRepo)
}

func (m Mother) Enqueuer() services.Enqueuer {
	return services.NewEnqueuer(m.Queue(), uuid.NewV4, m.MessagesRepo())
}

func (m Mother) TemplatesLoader() postal.TemplatesLoader {
	finder := m.TemplateFinder()
	database := m.Database()
	clientsRepo, kindsRepo := m.Repos()
	templatesRepo := m.TemplatesRepo()

	return postal.NewTemplatesLoader(finder, database, clientsRepo, kindsRepo, templatesRepo)
}

func (m Mother) UserLoader() postal.UserLoader {
	uaaClient := m.UAAClient()

	return postal.NewUserLoader(uaaClient)
}

func (m Mother) TokenLoader() postal.TokenLoader {
	uaaClient := m.UAAClient()

	return postal.NewTokenLoader(uaaClient)
}

func (m Mother) MailClient() *mail.Client {
	env := NewEnvironment()
	mailConfig := mail.Config{
		User:           env.SMTPUser,
		Pass:           env.SMTPPass,
		Host:           env.SMTPHost,
		Port:           env.SMTPPort,
		Secret:         env.SMTPCRAMMD5Secret,
		TestMode:       env.TestMode,
		SkipVerifySSL:  !env.VerifySSL,
		DisableTLS:     !env.SMTPTLS,
		LoggingEnabled: env.SMTPLoggingEnabled,
	}

	switch env.SMTPAuthMechanism {
	case SMTPAuthNone:
		mailConfig.AuthMechanism = mail.AuthNone
	case SMTPAuthPlain:
		mailConfig.AuthMechanism = mail.AuthPlain
	case SMTPAuthCRAMMD5:
		mailConfig.AuthMechanism = mail.AuthCRAMMD5
	}

	return mail.NewClient(mailConfig)
}

func (m Mother) Repos() (models.ClientsRepo, models.KindsRepo) {
	return models.NewClientsRepo(), m.KindsRepo()
}

func (m Mother) Logger() lager.Logger {
	logger := lager.NewLogger("notifications")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))

	return logger
}

func (m Mother) Logging() web.RequestLogging {
	return web.NewRequestLogging(m.Logger())
}

func (m Mother) ErrorWriter() handlers.ErrorWriter {
	return handlers.NewErrorWriter()
}

func (m Mother) Authenticator(scopes ...string) web.Authenticator {
	return web.NewAuthenticator(UAAPublicKey, scopes...)
}

func (m Mother) Registrar() services.Registrar {
	clientsRepo, kindsRepo := m.Repos()
	return services.NewRegistrar(clientsRepo, kindsRepo)
}

func (m *Mother) SQLDatabase() *sql.DB {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.sqlDB != nil {
		return m.sqlDB
	}

	env := NewEnvironment()

	var err error
	m.sqlDB, err = sql.Open("mysql", env.DatabaseURL)
	if err != nil {
		panic(err)
	}

	if err := m.sqlDB.Ping(); err != nil {
		panic(err)
	}

	m.sqlDB.SetMaxOpenConns(env.DBMaxOpenConns)

	return m.sqlDB
}

func (m *Mother) Database() models.DatabaseInterface {
	env := NewEnvironment()

	database := models.NewDatabase(m.SQLDatabase(), models.Config{
		DefaultTemplatePath: path.Join(env.RootPath, "templates", "default.json"),
	})
	if env.DBLoggingEnabled {
		database.TraceOn("[DB]", log.New(os.Stdout, "", 0))
	}
	database.Setup()
	return database
}

func (m Mother) PreferencesFinder() *services.PreferencesFinder {
	return services.NewPreferencesFinder(models.NewPreferencesRepo(), m.GlobalUnsubscribesRepo())
}

func (m Mother) PreferenceUpdater() services.PreferenceUpdater {
	return services.NewPreferenceUpdater(m.GlobalUnsubscribesRepo(), m.UnsubscribesRepo(), m.KindsRepo())
}

func (m Mother) TemplateFinder() services.TemplateFinder {
	return services.NewTemplateFinder(m.TemplatesRepo())
}

func (m Mother) MessageFinder() services.MessageFinder {
	return services.NewMessageFinder(m.MessagesRepo())
}

func (m Mother) TemplateServiceObjects() (services.TemplateCreator, services.TemplateFinder, services.TemplateUpdater,
	services.TemplateDeleter, services.TemplateLister, services.TemplateAssigner, services.TemplateAssociationLister) {

	clientsRepo, kindsRepo := m.Repos()
	templatesRepo := m.TemplatesRepo()

	return services.NewTemplateCreator(templatesRepo),
		m.TemplateFinder(),
		services.NewTemplateUpdater(templatesRepo),
		services.NewTemplateDeleter(templatesRepo),
		services.NewTemplateLister(templatesRepo),
		services.NewTemplateAssigner(clientsRepo, kindsRepo, templatesRepo),
		services.NewTemplateAssociationLister(clientsRepo, kindsRepo, templatesRepo)
}

func (m Mother) KindsRepo() models.KindsRepo {
	return models.NewKindsRepo()
}

func (m Mother) UnsubscribesRepo() models.UnsubscribesRepo {
	return models.NewUnsubscribesRepo()
}

func (m Mother) GlobalUnsubscribesRepo() models.GlobalUnsubscribesRepo {
	return models.NewGlobalUnsubscribesRepo()
}

func (m Mother) TemplatesRepo() models.TemplatesRepo {
	return models.NewTemplatesRepo()
}

func (m Mother) MessagesRepo() models.MessagesRepo {
	return models.NewMessagesRepo()
}

func (m Mother) ReceiptsRepo() models.ReceiptsRepo {
	return models.NewReceiptsRepo()
}

func (m Mother) CORS() web.CORS {
	env := NewEnvironment()
	return web.NewCORS(env.CORSOrigin)
}
