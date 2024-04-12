package application

import (
	"errors"
	"log"
	"os"
	"path"
	"time"

	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/mail"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/uaa"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/web"
	"github.com/pivotal-cf-experimental/warrant"
	"github.com/pivotal-golang/lager"
)

const WorkerCount = 10

type Application struct {
	env        Environment
	logger     lager.Logger
	dbProvider *DBProvider
	migrator   Migrator
}

func New(env Environment, dbp *DBProvider) Application {
	databaseMigrator := models.DatabaseMigrator{}

	l := lager.NewLogger("notifications")
	l.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))

	return Application{
		env:        env,
		logger:     l,
		dbProvider: dbp,
		migrator:   NewMigrator(dbp, databaseMigrator, env.VCAPApplication.InstanceIndex == 0, env.ModelMigrationsPath, env.GobbleMigrationsPath, path.Join(env.RootPath, "templates", "default.json")),
	}
}

func (a Application) mailClient() *mail.Client {
	return mail.NewClient(mail.Config{
		User:              a.env.SMTPUser,
		Pass:              a.env.SMTPPass,
		Host:              a.env.SMTPHost,
		Port:              a.env.SMTPPort,
		Secret:            a.env.SMTPCRAMMD5Secret,
		TestMode:          a.env.TestMode,
		SkipVerifySSL:     !a.env.VerifySSL,
		DisableTLS:        !a.env.SMTPTLS,
		LoggingEnabled:    a.env.SMTPLoggingEnabled,
		SMTPAuthMechanism: a.env.SMTPAuthMechanism,
	})
}

func (a Application) Run() {

	a.VerifySMTPConfiguration()

	uaaClient := warrant.New(warrant.Config{
		Host:          a.env.UAAHost,
		SkipVerifySSL: !a.env.VerifySSL,
	})

	validator := uaa.NewTokenValidator(a.logger, &uaaClient.Tokens)

	if err := validator.LoadSigningKeys(); err != nil {
		a.logger.Fatal("uaa-get-token-key-errored", err)
	}

	a.migrator.Migrate()

	a.StartQueueGauge()
	a.StartWorkers(validator)
	a.StartMessageGC()
	a.StartKeyRefresher(validator)
	a.StartServer(a.logger, validator)
}

func (a Application) VerifySMTPConfiguration() {
	if a.env.TestMode {
		return
	}

	mc := a.mailClient()
	err := mc.Connect(a.logger)
	if err != nil {
		a.logger.Fatal("smtp-connect-errored", err)
	}

	err = mc.Hello()
	if err != nil {
		a.logger.Fatal("smtp-hello-errored", err)
	}

	startTLSSupported, _ := mc.Extension("STARTTLS")

	mc.Quit()

	if !startTLSSupported && a.env.SMTPTLS {
		a.logger.Fatal("smtp-config-mismatch", errors.New(`SMTP TLS configuration mismatch: Configured to use TLS over SMTP, but the mail server does not support the "STARTTLS" extension.`))
	}

	if startTLSSupported && !a.env.SMTPTLS {
		a.logger.Fatal("smtp-config-mismatch", errors.New(`SMTP TLS configuration mismatch: Not configured to use TLS over SMTP, but the mail server does support the "STARTTLS" extension.`))
	}
}

func (a Application) StartQueueGauge() {
	if a.env.VCAPApplication.InstanceIndex != 0 {
		return
	}

	queueGauge := gobble.NewQueueGauge(a.dbProvider.Queue(), time.Tick(time.Minute))
	go queueGauge.Run()
}

func (a Application) StartKeyRefresher(validator *uaa.TokenValidator) {
	duration := time.Duration(a.env.UAAKeyRefreshInterval) * time.Millisecond

	t := time.NewTimer(duration)

	go func() {
		for {
			select {
			case <-t.C:
				validator.LoadSigningKeys()
				t.Reset(duration)
				break
			}
		}
	}()
}

func (a Application) StartWorkers(validator *uaa.TokenValidator) {
	postal.Boot(a.mailClient, a.dbProvider.sqlDB, postal.Config{
		UAAClientID:          a.env.UAAClientID,
		UAAClientSecret:      a.env.UAAClientSecret,
		UAATokenValidator:    validator,
		UAAHost:              a.env.UAAHost,
		VerifySSL:            a.env.VerifySSL,
		InstanceIndex:        a.env.VCAPApplication.InstanceIndex,
		WorkerCount:          WorkerCount,
		RootPath:             a.env.RootPath,
		EncryptionKey:        a.env.EncryptionKey,
		DBLoggingEnabled:     a.env.DBLoggingEnabled,
		Sender:               a.env.Sender,
		Domain:               a.env.Domain,
		QueueWaitMaxDuration: a.env.GobbleWaitMaxDuration,
		MaxQueueLength:       a.env.GobbleMaxQueueLength,
		MaxRetries:           a.env.MaxRetries,
		CCHost:               a.env.CCHost,
	})
}

func (a Application) StartMessageGC() {
	messageLifetime := 24 * time.Hour
	db := a.dbProvider.Database()
	messagesRepo := a.dbProvider.MessagesRepo()
	pollingInterval := 1 * time.Hour

	logger := log.New(os.Stdout, "", 0)
	messageGC := postal.NewMessageGC(messageLifetime, db, messagesRepo, pollingInterval, logger)
	messageGC.Run()
}

func (a Application) StartServer(logger lager.Logger, validator *uaa.TokenValidator) {
	web.NewServer().Run(web.Config{
		DBLoggingEnabled:     a.env.DBLoggingEnabled,
		SkipVerifySSL:        !a.env.VerifySSL,
		Port:                 a.env.Port,
		Logger:               logger,
		CORSOrigin:           a.env.CORSOrigin,
		SQLDB:                a.dbProvider.sqlDB,
		Queue:                a.dbProvider.Queue(),
		QueueWaitMaxDuration: a.env.GobbleWaitMaxDuration,
		MaxQueueLength:       a.env.GobbleMaxQueueLength,

		UAATokenValidator: validator,
		UAAHost:           a.env.UAAHost,
		UAAClientID:       a.env.UAAClientID,
		UAAClientSecret:   a.env.UAAClientSecret,
		DefaultUAAScopes:  a.env.DefaultUAAScopes,
		CCHost:            a.env.CCHost,
	})
}

// This is a hack to get the logs output to the loggregator before the process exits
func (a Application) Crash() {
	err := recover()
	switch err.(type) {
	case error:
		time.Sleep(5 * time.Second)
		a.logger.Fatal("crash", err.(error))
	case nil:
		return
	default:
		time.Sleep(5 * time.Second)
		a.logger.Fatal("crash", nil)
	}
}
