package postal

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/postal/common"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/pivotal-golang/lager"
	"github.com/rcrowley/go-metrics"
)

type DeliveryJobProcessor interface {
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
	DeliveryJobProcessor   DeliveryJobProcessor
	V2DeliveryJobProcessor v2DeliveryJobProcessor
	logger                 lager.Logger
	database               db.DatabaseInterface
	campaignJobProcessor   campaignJobProcessor
	deliveryFailureHandler deliveryFailureHandler
	messageStatusUpdater   messageStatusUpdater
}

func NewDeliveryWorker(v1DeliveryJobProcessor DeliveryJobProcessor, config DeliveryWorkerConfig) DeliveryWorker {
	worker := DeliveryWorker{
		DeliveryJobProcessor:   v1DeliveryJobProcessor,
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
		metrics.GetOrRegisterCounter("notifications.worker.panic.json", nil).Inc(1)

		worker.deliveryFailureHandler.Handle(job, worker.logger)
		return
	}

	worker.DeliveryJobProcessor.Process(job, worker.logger)
}
