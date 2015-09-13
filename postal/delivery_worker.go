package postal

import (
	"fmt"
	"time"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/pivotal-golang/lager"
)

type tokenLoader interface {
	Load(string) (string, error)
}

type process interface {
	Deliver(job *gobble.Job, logger lager.Logger) error
}

type workflow interface {
	Deliver(delivery Delivery, logger lager.Logger) error
}

type strategyDeterminer interface {
	Determine(conn services.ConnectionInterface, uaaHost string, job gobble.Job) error
}

type messageStatusUpdater interface {
	Update(conn db.ConnectionInterface, messageID, messageStatus, campaignID string, logger lager.Logger)
}

type deliveryFailureHandler interface {
	Handle(job Retryable, logger lager.Logger)
}

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

type DeliveryWorkerConfig struct {
	ID                     int
	UAAHost                string
	Logger                 lager.Logger
	Queue                  gobble.QueueInterface
	DBTrace                bool
	Database               db.DatabaseInterface
	StrategyDeterminer     strategyDeterminer
	DeliveryFailureHandler deliveryFailureHandler
	MessageStatusUpdater   messageStatusUpdater
}

type DeliveryWorker struct {
	gobble.Worker

	uaaHost                string
	V1Process              process
	V2Workflow             workflow
	logger                 lager.Logger
	database               db.DatabaseInterface
	strategyDeterminer     strategyDeterminer
	deliveryFailureHandler deliveryFailureHandler
	messageStatusUpdater   messageStatusUpdater
}

func NewDeliveryWorker(v1workflow process, v2workflow workflow, config DeliveryWorkerConfig) DeliveryWorker {
	logger := config.Logger.Session("worker", lager.Data{"worker_id": config.ID})

	worker := DeliveryWorker{
		V1Process:              v1workflow,
		V2Workflow:             v2workflow,
		uaaHost:                config.UAAHost,
		logger:                 logger,
		database:               config.Database,
		strategyDeterminer:     config.StrategyDeterminer,
		deliveryFailureHandler: config.DeliveryFailureHandler,
		messageStatusUpdater:   config.MessageStatusUpdater,
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
		var delivery Delivery
		job.Unmarshal(&delivery)

		err = worker.V2Workflow.Deliver(delivery, worker.logger)
		if err != nil {
			worker.deliveryFailureHandler.Handle(job, worker.logger)
			status := StatusFailed
			if job.ShouldRetry {
				status = StatusRetry
			}

			worker.messageStatusUpdater.Update(worker.database.Connection(), delivery.MessageID, status, delivery.CampaignID, worker.logger)
		}
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
