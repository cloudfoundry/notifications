package application

import (
	"log"
	"os"
	"time"

	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/web"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"
	"github.com/ryanmoran/viron"
)

const WorkerCount = 10

type Application struct {
	env      Environment
	mother   *Mother
	migrator Migrator
}

func NewApplication() Application {
	mother := NewMother()
	env := NewEnvironment()

	return Application{
		env:      env,
		mother:   mother,
		migrator: NewMigrator(mother, env.VCAPApplication.InstanceIndex == 0),
	}
}

func BootLogger() *log.Logger {
	return log.New(os.Stdout, "[BOOT] ", 0)
}

func (app Application) Boot(logger *log.Logger) {
	app.PrintConfiguration(logger)
	app.ConfigureSMTP(logger)
	app.RetrieveUAAPublicKey(logger)
	app.migrator.Migrate()
	app.EnableDBLogging()
	app.UnlockJobs()
	app.StartWorkers()
	app.StartMessageGC()
	app.StartServer()
}

func (app Application) PrintConfiguration(logger *log.Logger) {
	logger.Println("Booting with configuration:")

	viron.Print(app.env, logger)
}

func (app Application) ConfigureSMTP(logger *log.Logger) {
	if app.env.TestMode {
		return
	}

	mailClient := app.mother.MailClient()
	err := mailClient.Connect()
	if err != nil {
		logger.Panicln(err)
	}

	err = mailClient.Hello()
	if err != nil {
		logger.Panicln(err)
	}

	startTLSSupported, _ := mailClient.Extension("STARTTLS")

	mailClient.Quit()

	if !startTLSSupported && app.env.SMTPTLS {
		logger.Panicln(`SMTP TLS configuration mismatch: Configured to use TLS over SMTP, but the mail server does not support the "STARTTLS" extension.`)
	}

	if startTLSSupported && !app.env.SMTPTLS {
		logger.Panicln(`SMTP TLS configuration mismatch: Not configured to use TLS over SMTP, but the mail server does support the "STARTTLS" extension.`)
	}
}

func (app Application) RetrieveUAAPublicKey(logger *log.Logger) {
	uaaClient := app.mother.UAAClient()

	key, err := uaa.GetTokenKey(*uaaClient)
	if err != nil {
		logger.Panicln(err)
	}

	UAAPublicKey = key
	log.Printf("UAA Public Key: %s", UAAPublicKey)
}

func (app Application) UnlockJobs() {
	app.mother.Queue().Unlock()
}

func (app Application) EnableDBLogging() {
	if app.env.DBLoggingEnabled {
		app.mother.Database().TraceOn("[DB]", log.New(os.Stdout, "", 0))
	}
}

func (app Application) StartWorkers() {
	logger := log.New(os.Stdout, "", 0)

	WorkerGenerator{
		InstanceIndex: app.env.VCAPApplication.InstanceIndex,
		Count:         WorkerCount,
	}.Work(func(i int) Worker {
		worker := postal.NewDeliveryWorker(i, logger, app.mother.MailClient(), app.mother.Queue(),
			app.mother.GlobalUnsubscribesRepo(), app.mother.UnsubscribesRepo(), app.mother.KindsRepo(), app.mother.MessagesRepo(),
			app.mother.Database(), app.env.Sender, app.env.EncryptionKey, app.mother.UserLoader(), app.mother.TemplatesLoader(), app.mother.ReceiptsRepo(), app.mother.TokenLoader())
		return &worker
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

func (app Application) StartServer() {
	web.NewServer().Run(app.env.Port, app.mother)
}

// This is a hack to get the logs output to the loggregator before the process exits
func (app Application) Crash(logger *log.Logger) {
	err := recover()
	if err != nil {
		time.Sleep(5 * time.Second)
		logger.Panicln(err)
	}
}
