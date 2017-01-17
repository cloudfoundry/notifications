package application

import (
	"errors"
	"log"
	"os"
	"path"
	"time"

	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/uaa"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/web"
	"github.com/pivotal-golang/lager"
	"github.com/ryanmoran/viron"
)

const WorkerCount = 10

type Application struct {
	env      Environment
	mother   *Mother
	migrator Migrator
}

func NewApplication(env Environment, mother *Mother) Application {
	databaseMigrator := models.DatabaseMigrator{}
	return Application{
		env:      env,
		mother:   mother,
		migrator: NewMigrator(mother, databaseMigrator, env.VCAPApplication.InstanceIndex == 0, env.ModelMigrationsPath, env.GobbleMigrationsPath, path.Join(env.RootPath, "templates", "default.json")),
	}
}

func (app Application) Boot() {
	session := app.mother.Logger().Session("boot")

	viron.Print(app.env, vironCompatibleLogger{session})

	app.ConfigureSMTP(session)
	uaaKey := app.RetrieveUAAPublicKey(session)

	app.migrator.Migrate()

	app.StartQueueGauge()
	app.StartWorkers(uaaKey)
	app.StartMessageGC()
	app.StartServer(session, uaaKey)
}

func (app Application) ConfigureSMTP(logger lager.Logger) {
	if app.env.TestMode {
		return
	}

	mailClient := app.mother.MailClient()
	err := mailClient.Connect(logger)
	if err != nil {
		logger.Fatal("smtp-connect-errored", err)
	}

	err = mailClient.Hello()
	if err != nil {
		logger.Fatal("smtp-hello-errored", err)
	}

	startTLSSupported, _ := mailClient.Extension("STARTTLS")

	mailClient.Quit()

	if !startTLSSupported && app.env.SMTPTLS {
		logger.Fatal("smtp-config-mismatch", errors.New(`SMTP TLS configuration mismatch: Configured to use TLS over SMTP, but the mail server does not support the "STARTTLS" extension.`))
	}

	if startTLSSupported && !app.env.SMTPTLS {
		logger.Fatal("smtp-config-mismatch", errors.New(`SMTP TLS configuration mismatch: Not configured to use TLS over SMTP, but the mail server does support the "STARTTLS" extension.`))
	}
}

func (app Application) RetrieveUAAPublicKey(logger lager.Logger) string {
	zonedUAAClient := uaa.NewZonedUAAClient(app.env.UAAClientID, app.env.UAAClientSecret, app.env.VerifySSL, "")

	key, err := zonedUAAClient.GetTokenKey(app.env.UAAHost)
	if err != nil {
		logger.Fatal("uaa-get-token-key-errored", err)
	}

	logger.Info("uaa-public-key", lager.Data{
		"key": key,
	})

	return key
}

func (app Application) StartQueueGauge() {
	if app.env.VCAPApplication.InstanceIndex != 0 {
		return
	}

	queueGauge := metrics.NewQueueGauge(app.mother.Queue(), metrics.DefaultLogger, time.Tick(1*time.Second))
	go queueGauge.Run()
}

func (app Application) StartWorkers(uaaPublicKey string) {
	postal.Boot(app.mother, postal.Config{
		UAAClientID:          app.env.UAAClientID,
		UAAClientSecret:      app.env.UAAClientSecret,
		UAAPublicKey:         uaaPublicKey,
		UAAHost:              app.env.UAAHost,
		VerifySSL:            app.env.VerifySSL,
		InstanceIndex:        app.env.VCAPApplication.InstanceIndex,
		WorkerCount:          WorkerCount,
		EncryptionKey:        app.env.EncryptionKey,
		DBLoggingEnabled:     app.env.DBLoggingEnabled,
		Sender:               app.env.Sender,
		Domain:               app.env.Domain,
		QueueWaitMaxDuration: app.env.GobbleWaitMaxDuration,
		CCHost:               app.env.CCHost,
	})
}

func (app Application) StartMessageGC() {
	messageLifetime := 24 * time.Hour
	db := app.mother.Database()
	messagesRepo := app.mother.MessagesRepo()
	pollingInterval := 1 * time.Hour

	logger := log.New(os.Stdout, "", 0)
	messageGC := postal.NewMessageGC(messageLifetime, db, messagesRepo, pollingInterval, logger)
	messageGC.Run()
}

func (app Application) StartServer(logger lager.Logger, uaaPublicKey string) {
	web.NewServer().Run(app.mother, web.Config{
		DBLoggingEnabled:     app.env.DBLoggingEnabled,
		SkipVerifySSL:        !app.env.VerifySSL,
		Port:                 app.env.Port,
		Logger:               logger,
		CORSOrigin:           app.env.CORSOrigin,
		SQLDB:                app.mother.SQLDatabase(),
		QueueWaitMaxDuration: app.env.GobbleWaitMaxDuration,

		UAAPublicKey:     uaaPublicKey,
		UAAHost:          app.env.UAAHost,
		UAAClientID:      app.env.UAAClientID,
		UAAClientSecret:  app.env.UAAClientSecret,
		DefaultUAAScopes: app.env.DefaultUAAScopes,
		CCHost:           app.env.CCHost,
	})
}

// This is a hack to get the logs output to the loggregator before the process exits
func (app Application) Crash() {
	logger := app.mother.Logger()

	err := recover()
	switch err.(type) {
	case error:
		time.Sleep(5 * time.Second)
		logger.Fatal("crash", err.(error))
	case nil:
		return
	default:
		time.Sleep(5 * time.Second)
		logger.Fatal("crash", nil)
	}
}

type vironCompatibleLogger struct {
	logger lager.Logger
}

func (l vironCompatibleLogger) Printf(format string, v ...interface{}) {
	if len(v) == 2 {
		key, ok := v[0].(string)
		value := v[1]
		if ok {
			l.logger.Info("viron", lager.Data{key: value})
		}
	}
}
