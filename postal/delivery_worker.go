package postal

import (
	"log"
	"math"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/mail"
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/pivotal-golang/conceal"
)

type Delivery struct {
	Options      Options
	UserGUID     string
	Email        string
	Space        cf.CloudControllerSpace
	Organization cf.CloudControllerOrganization
	ClientID     string
	MessageID    string
	Scope        string
}

type MessagesRepoInterface interface {
	Upsert(models.ConnectionInterface, models.Message) (models.Message, error)
}

type DeliveryWorker struct {
	logger                 *log.Logger
	mailClient             mail.ClientInterface
	globalUnsubscribesRepo models.GlobalUnsubscribesRepoInterface
	unsubscribesRepo       models.UnsubscribesRepoInterface
	kindsRepo              models.KindsRepoInterface
	userLoader             UserLoaderInterface
	templatesLoader        TemplatesLoaderInterface
	tokenLoader            TokenLoaderInterface
	messagesRepo           MessagesRepoInterface
	receiptsRepo           models.ReceiptsRepoInterface
	database               models.DatabaseInterface
	sender                 string
	encryptionKey          []byte
	gobble.Worker
}

func NewDeliveryWorker(id int, logger *log.Logger, mailClient mail.ClientInterface, queue gobble.QueueInterface,
	globalUnsubscribesRepo models.GlobalUnsubscribesRepoInterface, unsubscribesRepo models.UnsubscribesRepoInterface,
	kindsRepo models.KindsRepoInterface, messagesRepo MessagesRepoInterface,
	database models.DatabaseInterface, sender string, encryptionKey []byte, userLoader UserLoaderInterface,
	templatesLoader TemplatesLoaderInterface, receiptsRepo models.ReceiptsRepoInterface, tokenLoader TokenLoaderInterface) DeliveryWorker {

	worker := DeliveryWorker{
		logger:                 logger,
		mailClient:             mailClient,
		globalUnsubscribesRepo: globalUnsubscribesRepo,
		unsubscribesRepo:       unsubscribesRepo,
		kindsRepo:              kindsRepo,
		messagesRepo:           messagesRepo,
		database:               database,
		sender:                 sender,
		encryptionKey:          encryptionKey,
		userLoader:             userLoader,
		tokenLoader:            tokenLoader,
		templatesLoader:        templatesLoader,
		receiptsRepo:           receiptsRepo,
	}
	worker.Worker = gobble.NewWorker(id, queue, worker.Deliver)

	return worker
}

func (worker DeliveryWorker) Deliver(job *gobble.Job) {
	var delivery Delivery

	err := job.Unmarshal(&delivery)
	if err != nil {
		metrics.NewMetric("counter", map[string]interface{}{
			"name": "notifications.worker.panic.json",
		}).Log()

		worker.retry(job)
		return
	}

	err = worker.receiptsRepo.CreateReceipts(worker.database.Connection(), []string{delivery.UserGUID}, delivery.ClientID, delivery.Options.KindID)
	if err != nil {
		worker.retry(job)
		return
	}

	if delivery.Email == "" {
		token, err := worker.tokenLoader.Load()
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
	message, err := worker.pack(delivery)
	if err != nil {
		worker.logger.Printf("Not delivering because template failed to pack")
		worker.updateMessageStatus(delivery.MessageID, StatusFailed)
		return StatusFailed
	}

	status := worker.sendMail(message)
	worker.updateMessageStatus(delivery.MessageID, status)

	return status
}

func (worker DeliveryWorker) updateMessageStatus(messageID, status string) {
	_, err := worker.messagesRepo.Upsert(worker.database.Connection(), models.Message{ID: messageID, Status: status})
	if err != nil {
		worker.logger.Printf("Failed to upsert status '%s' of notification %s. Error: %s", status, messageID, err.Error())
	}
}

func (worker DeliveryWorker) retry(job *gobble.Job) {
	if job.RetryCount < 10 {
		duration := time.Duration(int64(math.Pow(2, float64(job.RetryCount))))
		job.Retry(duration * time.Minute)
		layout := "Jan 2, 2006 at 3:04pm (MST)"
		worker.logger.Printf("Message failed to send, retrying at: %s", job.ActiveAt.Format(layout))
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
		worker.logger.Printf("Not delivering because %s has unsubscribed", delivery.Email)
		return false
	}

	_, err = worker.unsubscribesRepo.Find(conn, delivery.ClientID, delivery.Options.KindID, delivery.UserGUID)
	if err != nil {
		if _, ok := err.(models.RecordNotFoundError); ok {
			if delivery.Email == "" {
				worker.logger.Printf("Not delivering because recipient has no email addresses")
				return false
			}

			if !strings.Contains(delivery.Email, "@") {
				worker.logger.Printf("Not delivering because recipient's email address is invalid")
				return false
			}

			return true
		}

		worker.logger.Printf("Not delivering because: %+v", err)
		return false
	}

	worker.logger.Printf("Not delivering because %s has unsubscribed", delivery.Email)
	return false
}

func (worker DeliveryWorker) isCritical(conn models.ConnectionInterface, kindID, clientID string) bool {
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

	context := NewMessageContext(delivery, worker.sender, cloak, templates)
	packager := NewPackager()

	message, err = packager.Pack(context)
	if err != nil {
		return message, err
	}

	return message, nil
}

func (worker DeliveryWorker) sendMail(message mail.Message) string {
	err := worker.mailClient.Connect()
	if err != nil {
		worker.logger.Printf("Error Establishing SMTP Connection: %s", err.Error())
		return StatusUnavailable
	}

	worker.logger.Printf("Attempting to deliver message to %s", message.To)
	err = worker.mailClient.Send(message)
	if err != nil {
		worker.logger.Printf("Failed to deliver message due to SMTP error: %s", err.Error())
		return StatusFailed
	}

	worker.logger.Printf("Message was successfully sent to %s", message.To)

	return StatusDelivered
}
