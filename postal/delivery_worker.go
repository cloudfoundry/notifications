package postal

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/postal/common"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/pivotal-golang/lager"
)

type v1DeliveryJobProcessor interface {
	Process(job *gobble.Job, logger lager.Logger) error
}

type v2DeliveryJobProcessor interface {
	Process(delivery common.Delivery, logger lager.Logger) error
}

type campaignJobProcessor interface {
	Process(conn services.ConnectionInterface, uaaHost string, job gobble.Job, logger lager.Logger) error
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
	V1DeliveryJobProcessor v1DeliveryJobProcessor
	V2DeliveryJobProcessor v2DeliveryJobProcessor
	logger                 lager.Logger
	database               db.DatabaseInterface
	campaignJobProcessor   campaignJobProcessor
	deliveryFailureHandler deliveryFailureHandler
	messageStatusUpdater   messageStatusUpdater
}

func NewDeliveryWorker(v1DeliveryJobProcessor v1DeliveryJobProcessor, v2DeliveryJobProcessor v2DeliveryJobProcessor, config DeliveryWorkerConfig) DeliveryWorker {
	worker := DeliveryWorker{
		V1DeliveryJobProcessor: v1DeliveryJobProcessor,
		V2DeliveryJobProcessor: v2DeliveryJobProcessor,
		uaaHost:                config.UAAHost,
		logger:                 config.Logger,
		database:               config.Database,
		campaignJobProcessor:   config.CampaignJobProcessor,
		deliveryFailureHandler: config.DeliveryFailureHandler,
		messageStatusUpdater:   config.MessageStatusUpdater,
	}
	ticker := gobble.NewTicker(time.NewTicker, 30*time.Second)
	heartbeater := gobble.NewHeartbeater(config.Queue, ticker)
	worker.Worker = gobble.NewWorker(config.ID, config.Queue, worker.Deliver, heartbeater)

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
		err := worker.campaignJobProcessor.Process(worker.database.Connection(), worker.uaaHost, *job, worker.logger)
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
		worker.V1DeliveryJobProcessor.Process(job, worker.logger)
	}
}
