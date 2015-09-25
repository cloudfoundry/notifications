package application

import (
	"crypto/rand"
	"database/sql"
	"log"
	"os"
	"path"
	"sync"
	"time"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/mail"
	"github.com/cloudfoundry-incubator/notifications/uaa"
	v1models "github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	v2models "github.com/cloudfoundry-incubator/notifications/v2/models"
	"github.com/cloudfoundry-incubator/notifications/v2/queue"
	"github.com/cloudfoundry-incubator/notifications/v2/util"
	"github.com/pivotal-golang/lager"
)

type Mother struct {
	sqlDB *sql.DB
	mutex sync.Mutex
	env   Environment
}

func NewMother(env Environment) *Mother {
	return &Mother{
		env: env,
	}
}

func (m *Mother) GobbleDatabase() gobble.DatabaseInterface {
	return gobble.NewDatabase(m.SQLDatabase())
}

func (m *Mother) Queue() gobble.QueueInterface {
	return gobble.NewQueue(m.GobbleDatabase(), gobble.Config{
		WaitMaxDuration: time.Duration(m.env.GobbleWaitMaxDuration) * time.Millisecond,
	})
}

func (m *Mother) V2Enqueuer() queue.JobEnqueuer {
	return queue.NewJobEnqueuer(m.Queue(), v2models.NewMessagesRepository(util.NewClock(), v2models.NewIDGenerator(rand.Reader).Generate))
}

func (m *Mother) UserStrategy() services.UserStrategy {
	return services.NewUserStrategy(m.Enqueuer(), m.V2Enqueuer())
}

func (m *Mother) SpaceStrategy() services.SpaceStrategy {
	uaaClient := uaa.NewZonedUAAClient(m.env.UAAClientID, m.env.UAAClientSecret, m.env.VerifySSL, UAAPublicKey)
	cloudController := cf.NewCloudController(m.env.CCHost, !m.env.VerifySSL)

	tokenLoader := uaa.NewTokenLoader(uaaClient)
	spaceLoader := services.NewSpaceLoader(cloudController)
	organizationLoader := services.NewOrganizationLoader(cloudController)
	enqueuer := m.Enqueuer()
	findsUserIDs := services.NewFindsUserIDs(cloudController, uaaClient)

	return services.NewSpaceStrategy(tokenLoader, spaceLoader, organizationLoader, findsUserIDs, enqueuer, m.V2Enqueuer())
}

func (m *Mother) OrganizationStrategy() services.OrganizationStrategy {
	cloudController := cf.NewCloudController(m.env.CCHost, !m.env.VerifySSL)

	uaaClient := uaa.NewZonedUAAClient(m.env.UAAClientID, m.env.UAAClientSecret, m.env.VerifySSL, UAAPublicKey)
	tokenLoader := uaa.NewTokenLoader(uaaClient)
	organizationLoader := services.NewOrganizationLoader(cloudController)
	findsUserIDs := services.NewFindsUserIDs(cloudController, uaaClient)
	enqueuer := m.Enqueuer()

	return services.NewOrganizationStrategy(tokenLoader, organizationLoader, findsUserIDs, enqueuer, m.V2Enqueuer())
}

func (m *Mother) EveryoneStrategy() services.EveryoneStrategy {
	uaaClient := uaa.NewZonedUAAClient(m.env.UAAClientID, m.env.UAAClientSecret, m.env.VerifySSL, UAAPublicKey)
	tokenLoader := uaa.NewTokenLoader(uaaClient)
	allUsers := services.NewAllUsers(uaaClient)
	enqueuer := m.Enqueuer()

	return services.NewEveryoneStrategy(tokenLoader, allUsers, enqueuer, m.V2Enqueuer())
}

func (m *Mother) UAAScopeStrategy() services.UAAScopeStrategy {
	uaaClient := uaa.NewZonedUAAClient(m.env.UAAClientID, m.env.UAAClientSecret, m.env.VerifySSL, UAAPublicKey)
	cloudController := cf.NewCloudController(m.env.CCHost, !m.env.VerifySSL)

	tokenLoader := uaa.NewTokenLoader(uaaClient)
	findsUserIDs := services.NewFindsUserIDs(cloudController, uaaClient)
	enqueuer := m.Enqueuer()

	return services.NewUAAScopeStrategy(tokenLoader, findsUserIDs, enqueuer, m.V2Enqueuer(), m.env.DefaultUAAScopes)
}

func (m *Mother) EmailStrategy() services.EmailStrategy {
	return services.NewEmailStrategy(m.Enqueuer(), m.V2Enqueuer())
}

func (m *Mother) Enqueuer() services.Enqueuer {
	return services.NewEnqueuer(m.Queue(), m.MessagesRepo())
}

func (m *Mother) MailClient() *mail.Client {
	mailConfig := mail.Config{
		User:           m.env.SMTPUser,
		Pass:           m.env.SMTPPass,
		Host:           m.env.SMTPHost,
		Port:           m.env.SMTPPort,
		Secret:         m.env.SMTPCRAMMD5Secret,
		TestMode:       m.env.TestMode,
		SkipVerifySSL:  !m.env.VerifySSL,
		DisableTLS:     !m.env.SMTPTLS,
		LoggingEnabled: m.env.SMTPLoggingEnabled,
	}

	switch m.env.SMTPAuthMechanism {
	case SMTPAuthNone:
		mailConfig.AuthMechanism = mail.AuthNone
	case SMTPAuthPlain:
		mailConfig.AuthMechanism = mail.AuthPlain
	case SMTPAuthCRAMMD5:
		mailConfig.AuthMechanism = mail.AuthCRAMMD5
	}

	return mail.NewClient(mailConfig)
}

func (m *Mother) Logger() lager.Logger {
	logger := lager.NewLogger("notifications")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))

	return logger
}

func (m *Mother) SQLDatabase() *sql.DB {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.sqlDB != nil {
		return m.sqlDB
	}

	var err error
	m.sqlDB, err = sql.Open("mysql", m.env.DatabaseURL)
	if err != nil {
		panic(err)
	}

	if err := m.sqlDB.Ping(); err != nil {
		panic(err)
	}

	m.sqlDB.SetMaxOpenConns(m.env.DBMaxOpenConns)

	return m.sqlDB
}

func (m *Mother) Database() db.DatabaseInterface {
	database := v1models.NewDatabase(m.SQLDatabase(), v1models.Config{
		DefaultTemplatePath: path.Join(m.env.RootPath, "templates", "default.json"),
	})

	if m.env.DBLoggingEnabled {
		database.TraceOn("[DB]", log.New(os.Stdout, "", 0))
	}

	return database
}

func (m *Mother) KindsRepo() v1models.KindsRepo {
	return v1models.NewKindsRepo()
}

func (m *Mother) UnsubscribesRepo() v1models.UnsubscribesRepo {
	return v1models.NewUnsubscribesRepo()
}

func (m *Mother) GlobalUnsubscribesRepo() v1models.GlobalUnsubscribesRepo {
	return v1models.NewGlobalUnsubscribesRepo()
}

func (m *Mother) MessagesRepo() v1models.MessagesRepo {
	return v1models.NewMessagesRepo(v2models.NewIDGenerator(rand.Reader).Generate)
}

func (m *Mother) ReceiptsRepo() v1models.ReceiptsRepo {
	return v1models.NewReceiptsRepo()
}
