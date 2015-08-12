package postal

import (
	"fmt"
	"math"
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

	dbTrace       bool
	sender        string
	encryptionKey []byte
	identifier    int

	baseLogger lager.Logger
	logger     lager.Logger

	mailClient      mail.ClientInterface
	userLoader      UserLoaderInterface
	templatesLoader TemplatesLoaderInterface
	tokenLoader     TokenLoaderInterface
	database        db.DatabaseInterface

	messagesRepo           MessagesRepo
	receiptsRepo           ReceiptsRepo
	globalUnsubscribesRepo GlobalUnsubscribesRepo
	unsubscribesRepo       UnsubscribesRepo
	kindsRepo              KindsRepo
}

type DeliveryWorkerConfig struct {
	ID            int
	Sender        string
	Domain        string
	EncryptionKey []byte

	Logger                 lager.Logger
	MailClient             mail.ClientInterface
	Queue                  gobble.QueueInterface
	Database               db.DatabaseInterface
	DBTrace                bool
	GlobalUnsubscribesRepo GlobalUnsubscribesRepo
	UnsubscribesRepo       UnsubscribesRepo
	KindsRepo              KindsRepo
	MessagesRepo           MessagesRepo
	UserLoader             UserLoaderInterface
	TemplatesLoader        TemplatesLoaderInterface
	ReceiptsRepo           ReceiptsRepo
	TokenLoader            TokenLoaderInterface
}

type TokenLoaderInterface interface {
	Load(string) (string, error)
}

func NewDeliveryWorker(config DeliveryWorkerConfig) DeliveryWorker {

	logger := config.Logger.Session("worker", lager.Data{"worker_id": config.ID})

	worker := DeliveryWorker{
		identifier:             config.ID,
		baseLogger:             logger,
		logger:                 logger,
		mailClient:             config.MailClient,
		globalUnsubscribesRepo: config.GlobalUnsubscribesRepo,
		unsubscribesRepo:       config.UnsubscribesRepo,
		kindsRepo:              config.KindsRepo,
		messagesRepo:           config.MessagesRepo,
		database:               config.Database,
		dbTrace:                config.DBTrace,
		sender:                 config.Sender,
		encryptionKey:          config.EncryptionKey,
		userLoader:             config.UserLoader,
		tokenLoader:            config.TokenLoader,
		templatesLoader:        config.TemplatesLoader,
		receiptsRepo:           config.ReceiptsRepo,
	}
	worker.Worker = gobble.NewWorker(config.ID, config.Queue, worker.Deliver)

	return worker
}

func (worker DeliveryWorker) Deliver(job *gobble.Job) {
	var delivery Delivery

	err := job.Unmarshal(&delivery)
	if err != nil {
		metrics.NewMetric("counter", map[string]interface{}{
			"name": "notifications.worker.panic.json",
		}).Log()

		worker.retry("UNKNOWN", "UNKNOWN", job)
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
		worker.retry(delivery.MessageID, delivery.Email, job)
		return
	}

	if delivery.Email == "" {
		var token string

		token, err = worker.tokenLoader.Load(delivery.UAAHost)
		if err != nil {
			worker.retry(delivery.MessageID, delivery.Email, job)
			return
		}

		users, err := worker.userLoader.Load([]string{delivery.UserGUID}, token)
		if err != nil || len(users) < 1 {
			worker.retry(delivery.MessageID, delivery.Email, job)
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
			worker.retry(delivery.MessageID, delivery.Email, job)
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
	message, err := worker.pack(delivery)
	if err != nil {
		worker.logger.Info("template-pack-failed")
		worker.updateMessageStatus(delivery.MessageID, StatusFailed, delivery.Email)
		return StatusFailed
	}

	status := worker.sendMail(delivery.MessageID, message)
	worker.updateMessageStatus(delivery.MessageID, status, delivery.Email)

	return status
}

func (worker DeliveryWorker) updateMessageStatus(messageID, status, recipient string) {
	_, err := worker.messagesRepo.Upsert(worker.database.Connection(), models.Message{ID: messageID, Status: status})
	if err != nil {
		worker.logger.Error("failed-message-status-upsert", err, lager.Data{
			"status": status,
		})
	}
}

func (worker DeliveryWorker) retry(messageID, recipient string, job *gobble.Job) {
	if job.RetryCount < 10 {
		duration := time.Duration(int64(math.Pow(2, float64(job.RetryCount))))
		job.Retry(duration * time.Minute)

		worker.logger.Info("delivery-failed-retrying", lager.Data{
			"retry_count": job.RetryCount,
			"active_at":   job.ActiveAt.Format(time.RFC3339),
		})
	}

	metrics.NewMetric("counter", map[string]interface{}{
		"name": "notifications.worker.retry",
	}).Log()
}

func (worker DeliveryWorker) shouldDeliver(delivery Delivery) bool {
	conn := worker.database.Connection()
	if worker.isCritical(conn, delivery.Options.KindID, delivery.ClientID) {
		return true
	}

	globallyUnsubscribed, err := worker.globalUnsubscribesRepo.Get(conn, delivery.UserGUID)
	if err != nil || globallyUnsubscribed {
		worker.logger.Info("user-unsubscribed")
		worker.updateMessageStatus(delivery.MessageID, StatusUndeliverable, delivery.Email)
		return false
	}

	isUnsubscribed, err := worker.unsubscribesRepo.Get(conn, delivery.UserGUID, delivery.ClientID, delivery.Options.KindID)
	if err != nil || isUnsubscribed {
		worker.logger.Info("user-unsubscribed")
		worker.updateMessageStatus(delivery.MessageID, StatusUndeliverable, delivery.Email)
		return false
	}

	if delivery.Email == "" {
		worker.logger.Info("no-email-address-for-user")
		worker.updateMessageStatus(delivery.MessageID, StatusUndeliverable, delivery.Email)
		return false
	}

	if !strings.Contains(delivery.Email, "@") {
		worker.logger.Info("malformatted-email-address")
		worker.updateMessageStatus(delivery.MessageID, StatusUndeliverable, delivery.Email)
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

func (worker DeliveryWorker) pack(delivery Delivery) (mail.Message, error) {
	var message mail.Message

	cloak, err := conceal.NewCloak([]byte(worker.encryptionKey))
	if err != nil {
		panic(err)
	}

	templates, err := worker.templatesLoader.LoadTemplates(delivery.ClientID, delivery.Options.KindID)
	if err != nil {
		return message, err
	}

	context := NewMessageContext(delivery, worker.sender, "change-me", cloak, templates)
	packager := NewPackager()

	message, err = packager.Pack(context)
	if err != nil {
		return message, err
	}

	return message, nil
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
