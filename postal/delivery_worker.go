package postal

import (
	"fmt"
	"time"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/mail"
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
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
	CampaignID      string
}

type DeliveryWorker struct {
	gobble.Worker

	uaaHost                string
	V1Process              workflow
	V2Workflow             workflow
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

type workflow interface {
	Deliver(job *gobble.Job, logger lager.Logger) error
}

type StrategyDeterminerInterface interface {
	Determine(conn services.ConnectionInterface, uaaHost string, job gobble.Job) error
}

type messageStatusUpdaterInterface interface {
	Update(conn db.ConnectionInterface, messageID, messageStatus, campaignID string, logger lager.Logger)
}

type deliveryFailureHandlerInterface interface {
	Handle(job Retryable, logger lager.Logger)
}

func NewDeliveryWorker(v1workflow, v2workflow workflow, config DeliveryWorkerConfig) DeliveryWorker {
	logger := config.Logger.Session("worker", lager.Data{"worker_id": config.ID})

	worker := DeliveryWorker{
		V1Process:              v1workflow,
		V2Workflow:             v2workflow,
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
	var typedJob struct {
		JobType string
	}

	err := job.Unmarshal(&typedJob)
	if err != nil {
		metrics.NewMetric("counter", map[string]interface{}{
			"name": "notifications.worker.panic.json",
		}).Log()

		worker.deliveryFailureHandler.Handle(job, worker.logger)
		return
	}

	switch typedJob.JobType {
	case "campaign":
		err := worker.strategyDeterminer.Determine(worker.database.Connection(), worker.uaaHost, *job)
		if err != nil {
			worker.deliveryFailureHandler.Handle(job, worker.logger)
		}
	case "v2":
		worker.V2Workflow.Deliver(job, worker.logger)
	default:
		worker.V1Process.Deliver(job, worker.logger)
	}
}

type gorpCompatibleLogger struct {
	logger lager.Logger
}

func (g gorpCompatibleLogger) Printf(format string, v ...interface{}) {
	g.logger.Info("db", lager.Data{
		"statement": fmt.Sprintf(format, v...),
	})
}
