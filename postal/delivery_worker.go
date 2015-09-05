package postal

import (
	"fmt"
	"time"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/mail"
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/pivotal-golang/conceal"
	"github.com/pivotal-golang/lager"
)

type Delivery struct {
	MessageID       string
	Options         Options
	UserGUID        string
	Email           string
	Space           cf.CloudControllerSpace
	Organization    cf.CloudControllerOrganization
	ClientID        string
	UAAHost         string
	Scope           string
	VCAPRequestID   string
	RequestReceived time.Time
}

type DeliveryWorker struct {
	gobble.Worker

	uaaHost                string
	V1Process              V1Process
	logger                 lager.Logger
	database               db.DatabaseInterface
	strategyDeterminer     StrategyDeterminerInterface
	deliveryFailureHandler deliveryFailureHandlerInterface
}

type DeliveryWorkerConfig struct {
	ID            int
	Sender        string
	Domain        string
	EncryptionKey []byte
	UAAHost       string

	Logger                 lager.Logger
	MailClient             mail.ClientInterface
	Queue                  gobble.QueueInterface
	Database               db.DatabaseInterface
	DBTrace                bool
	GlobalUnsubscribesRepo GlobalUnsubscribesRepo
	UnsubscribesRepo       UnsubscribesRepo
	KindsRepo              KindsRepo
	UserLoader             UserLoaderInterface
	TemplatesLoader        TemplatesLoaderInterface
	ReceiptsRepo           ReceiptsRepo
	TokenLoader            TokenLoaderInterface
	StrategyDeterminer     StrategyDeterminerInterface
	MessageStatusUpdater   messageStatusUpdaterInterface
	DeliveryFailureHandler deliveryFailureHandlerInterface
}

type TokenLoaderInterface interface {
	Load(string) (string, error)
}

type StrategyDeterminerInterface interface {
	Determine(conn db.ConnectionInterface, uaaHost string, job gobble.Job) error
}

type messageStatusUpdaterInterface interface {
	Update(conn models.ConnectionInterface, messageID, messageStatus string, logger lager.Logger)
}

type deliveryFailureHandlerInterface interface {
	Handle(job Retryable, logger lager.Logger)
}

func NewDeliveryWorker(config DeliveryWorkerConfig) DeliveryWorker {
	logger := config.Logger.Session("worker", lager.Data{"worker_id": config.ID})

	cloak, err := conceal.NewCloak(config.EncryptionKey)
	if err != nil {
		panic(err)
	}

	worker := DeliveryWorker{
		V1Process: V1Process{
			dbTrace:                config.DBTrace,
			uaaHost:                config.UAAHost,
			sender:                 config.Sender,
			domain:                 config.Domain,
			packager:               NewPackager(config.TemplatesLoader, cloak),
			mailClient:             config.MailClient,
			database:               config.Database,
			tokenLoader:            config.TokenLoader,
			userLoader:             config.UserLoader,
			kindsRepo:              config.KindsRepo,
			receiptsRepo:           config.ReceiptsRepo,
			unsubscribesRepo:       config.UnsubscribesRepo,
			globalUnsubscribesRepo: config.GlobalUnsubscribesRepo,
			messageStatusUpdater:   config.MessageStatusUpdater,
			deliveryFailureHandler: config.DeliveryFailureHandler,
		},

		uaaHost:                config.UAAHost,
		logger:                 logger,
		database:               config.Database,
		strategyDeterminer:     config.StrategyDeterminer,
		deliveryFailureHandler: config.DeliveryFailureHandler,
	}
	worker.Worker = gobble.NewWorker(config.ID, config.Queue, worker.Deliver)

	return worker
}

func (worker DeliveryWorker) Deliver(job *gobble.Job) {
	var campaignJob struct {
		JobType string
	}

	err := job.Unmarshal(&campaignJob)
	if err != nil {
		metrics.NewMetric("counter", map[string]interface{}{
			"name": "notifications.worker.panic.json",
		}).Log()

		worker.deliveryFailureHandler.Handle(job, worker.logger)
		return
	}

	if campaignJob.JobType == "campaign" {
		worker.logger.Info("determining-strategy")
		if err := worker.strategyDeterminer.Determine(worker.database.Connection(), worker.uaaHost, *job); err != nil {
			worker.logger.Error("determining-strategy-failed", err)
			worker.deliveryFailureHandler.Handle(job, worker.logger)
		}

		return
	}

	worker.V1Process.Deliver(job, worker.logger)
}

type gorpCompatibleLogger struct {
	logger lager.Logger
}

func (g gorpCompatibleLogger) Printf(format string, v ...interface{}) {
	g.logger.Info("db", lager.Data{
		"statement": fmt.Sprintf(format, v...),
	})
}
