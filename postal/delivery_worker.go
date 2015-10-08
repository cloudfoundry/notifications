package postal

import (
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/postal/common"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/pivotal-golang/lager"
)

type process interface {
	Deliver(job *gobble.Job, logger lager.Logger) error
}

type v2DeliveryJobProcessor interface {
	Process(delivery common.Delivery, logger lager.Logger) error
}

type campaignJobProcessor interface {
	Process(conn services.ConnectionInterface, uaaHost string, job gobble.Job) error
}

type messageStatusUpdater interface {
	Update(conn db.ConnectionInterface, messageID, messageStatus, campaignID string, logger lager.Logger)
}

type deliveryFailureHandler interface {
	Handle(job common.Retryable, logger lager.Logger)
}

type DeliveryWorkerConfig struct {
	ID                     int
	UAAHost                string
	Logger                 lager.Logger
	Queue                  gobble.QueueInterface
	DBTrace                bool
	Database               db.DatabaseInterface
	CampaignJobProcessor   campaignJobProcessor
	DeliveryFailureHandler deliveryFailureHandler
	MessageStatusUpdater   messageStatusUpdater
}

type DeliveryWorker struct {
	gobble.Worker

	uaaHost                string
	V1Process              process
	V2DeliveryJobProcessor v2DeliveryJobProcessor
	logger                 lager.Logger
	database               db.DatabaseInterface
	campaignJobProcessor   campaignJobProcessor
	deliveryFailureHandler deliveryFailureHandler
	messageStatusUpdater   messageStatusUpdater
}

func NewDeliveryWorker(v1process process, v2DeliveryJobProcessor v2DeliveryJobProcessor, config DeliveryWorkerConfig) DeliveryWorker {
	logger := config.Logger.Session("worker", lager.Data{"worker_id": config.ID})

	worker := DeliveryWorker{
		V1Process:              v1process,
		V2DeliveryJobProcessor: v2DeliveryJobProcessor,
		uaaHost:                config.UAAHost,
		logger:                 logger,
		database:               config.Database,
		campaignJobProcessor:   config.CampaignJobProcessor,
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
		err := worker.campaignJobProcessor.Process(worker.database.Connection(), worker.uaaHost, *job)
		if err != nil {
			worker.deliveryFailureHandler.Handle(job, worker.logger)
		}
	case "v2":
		var delivery common.Delivery
		job.Unmarshal(&delivery)

		err = worker.V2DeliveryJobProcessor.Process(delivery, worker.logger)
		if err != nil {
			worker.deliveryFailureHandler.Handle(job, worker.logger)
			status := common.StatusFailed
			if job.ShouldRetry {
				status = common.StatusRetry
			}

			worker.messageStatusUpdater.Update(worker.database.Connection(), delivery.MessageID, status, delivery.CampaignID, worker.logger)
		}
	default:
		worker.V1Process.Deliver(job, worker.logger)
	}
}
