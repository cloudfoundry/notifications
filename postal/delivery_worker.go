package postal

import (
	"fmt"
	"strings"
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

	dbTrace    bool
	sender     string
	identifier int
	domain     string
	uaaHost    string

	baseLogger lager.Logger
	logger     lager.Logger

	packager           Packager
	mailClient         mail.ClientInterface
	userLoader         UserLoaderInterface
	templatesLoader    TemplatesLoaderInterface
	tokenLoader        TokenLoaderInterface
	strategyDeterminer StrategyDeterminerInterface
	database           db.DatabaseInterface

	receiptsRepo           ReceiptsRepo
	globalUnsubscribesRepo GlobalUnsubscribesRepo
	unsubscribesRepo       UnsubscribesRepo
	kindsRepo              KindsRepo
	messageStatusUpdater   messageStatusUpdaterInterface
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
		identifier:             config.ID,
		baseLogger:             logger,
		logger:                 logger,
		domain:                 config.Domain,
		uaaHost:                config.UAAHost,
		mailClient:             config.MailClient,
		globalUnsubscribesRepo: config.GlobalUnsubscribesRepo,
		unsubscribesRepo:       config.UnsubscribesRepo,
		kindsRepo:              config.KindsRepo,
		database:               config.Database,
		dbTrace:                config.DBTrace,
		sender:                 config.Sender,
		packager:               NewPackager(config.TemplatesLoader, cloak),
		userLoader:             config.UserLoader,
		tokenLoader:            config.TokenLoader,
		templatesLoader:        config.TemplatesLoader,
		receiptsRepo:           config.ReceiptsRepo,
		strategyDeterminer:     config.StrategyDeterminer,
		messageStatusUpdater:   config.MessageStatusUpdater,
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

		worker.retry(job)
		return
	}

	if campaignJob.JobType == "campaign" {
		worker.logger.Info("determining-strategy")
		if err := worker.strategyDeterminer.Determine(worker.database.Connection(), worker.uaaHost, *job); err != nil {
			worker.logger.Error("determining-strategy-failed", err)
		}
		worker.retry(job)
		return
	}

	var delivery Delivery
	err = job.Unmarshal(&delivery)
	if err != nil {
		metrics.NewMetric("counter", map[string]interface{}{
			"name": "notifications.worker.panic.json",
		}).Log()

		worker.retry(job)
		return
	}

	worker.logger = worker.logger.WithData(lager.Data{
		"message_id":      delivery.MessageID,
		"vcap_request_id": delivery.VCAPRequestID,
	})

	if worker.dbTrace {
		worker.database.TraceOn("", gorpCompatibleLogger{worker.logger})
	}

	err = worker.receiptsRepo.CreateReceipts(worker.database.Connection(), []string{delivery.UserGUID}, delivery.ClientID, delivery.Options.KindID)
	if err != nil {
		worker.retry(job)
		return
	}

	if delivery.Email == "" {
		var token string

		token, err = worker.tokenLoader.Load(worker.uaaHost)
		if err != nil {
			worker.retry(job)
			return
		}

		users, err := worker.userLoader.Load([]string{delivery.UserGUID}, token)
		if err != nil || len(users) < 1 {
			worker.retry(job)
			return
		}

		emails := users[delivery.UserGUID].Emails
		if len(emails) > 0 {
			delivery.Email = emails[0]
		}
	}

	worker.logger = worker.logger.WithData(lager.Data{
		"recipient": delivery.Email,
	})

	if worker.shouldDeliver(delivery) {
		status := worker.deliver(delivery)

		if status != StatusDelivered {
			worker.retry(job)
			return
		} else {
			metrics.NewMetric("counter", map[string]interface{}{
				"name": "notifications.worker.delivered",
			}).Log()
		}
	} else {
		metrics.NewMetric("counter", map[string]interface{}{
			"name": "notifications.worker.unsubscribed",
		}).Log()
	}
}

func (worker DeliveryWorker) deliver(delivery Delivery) string {
	context, err := worker.packager.PrepareContext(delivery, worker.sender, worker.domain)
	if err != nil {
		panic(err)
	}

	message, err := worker.packager.Pack(context)
	if err != nil {
		worker.logger.Info("template-pack-failed")
		worker.messageStatusUpdater.Update(worker.database.Connection(), delivery.MessageID, StatusFailed, worker.logger)
		return StatusFailed
	}

	status := worker.sendMail(delivery.MessageID, message)
	worker.messageStatusUpdater.Update(worker.database.Connection(), delivery.MessageID, status, worker.logger)

	return status
}

func (worker DeliveryWorker) retry(job *gobble.Job) {
	worker.deliveryFailureHandler.Handle(job, worker.logger)
}

func (worker DeliveryWorker) shouldDeliver(delivery Delivery) bool {
	conn := worker.database.Connection()
	if worker.isCritical(conn, delivery.Options.KindID, delivery.ClientID) {
		return true
	}

	globallyUnsubscribed, err := worker.globalUnsubscribesRepo.Get(conn, delivery.UserGUID)
	if err != nil || globallyUnsubscribed {
		worker.logger.Info("user-unsubscribed")
		worker.messageStatusUpdater.Update(worker.database.Connection(), delivery.MessageID, StatusUndeliverable, worker.logger)
		return false
	}

	isUnsubscribed, err := worker.unsubscribesRepo.Get(conn, delivery.UserGUID, delivery.ClientID, delivery.Options.KindID)
	if err != nil || isUnsubscribed {
		worker.logger.Info("user-unsubscribed")
		worker.messageStatusUpdater.Update(worker.database.Connection(), delivery.MessageID, StatusUndeliverable, worker.logger)
		return false
	}

	if delivery.Email == "" {
		worker.logger.Info("no-email-address-for-user")
		worker.messageStatusUpdater.Update(worker.database.Connection(), delivery.MessageID, StatusUndeliverable, worker.logger)
		return false
	}

	if !strings.Contains(delivery.Email, "@") {
		worker.logger.Info("malformatted-email-address")
		worker.messageStatusUpdater.Update(worker.database.Connection(), delivery.MessageID, StatusUndeliverable, worker.logger)
		return false
	}

	return true
}

func (worker DeliveryWorker) isCritical(conn db.ConnectionInterface, kindID, clientID string) bool {
	kind, err := worker.kindsRepo.Find(conn, kindID, clientID)
	if _, ok := err.(models.RecordNotFoundError); ok {
		return false
	}

	return kind.Critical
}

func (worker DeliveryWorker) sendMail(messageID string, message mail.Message) string {
	err := worker.mailClient.Connect(worker.logger)
	if err != nil {
		worker.logger.Error("smtp-connection-error", err)
		return StatusUnavailable
	}

	worker.logger.Info("delivery-start")

	err = worker.mailClient.Send(message, worker.logger)
	if err != nil {
		worker.logger.Error("delivery-failed-smtp-error", err)
		return StatusFailed
	}

	worker.logger.Info("message-sent")

	return StatusDelivered
}

type gorpCompatibleLogger struct {
	logger lager.Logger
}

func (g gorpCompatibleLogger) Printf(format string, v ...interface{}) {
	g.logger.Info("db", lager.Data{
		"statement": fmt.Sprintf(format, v...),
	})
}
