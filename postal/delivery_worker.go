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
	"github.com/pivotal-cf/uaa-sso-golang/uaa"
	"github.com/pivotal-golang/conceal"
)

type Delivery struct {
	User         uaa.User
	Options      Options
	UserGUID     string
	Space        cf.CloudControllerSpace
	Organization cf.CloudControllerOrganization
	ClientID     string
	Templates    Templates
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
	messagesRepo           MessagesRepoInterface
	database               models.DatabaseInterface
	sender                 string
	encryptionKey          []byte
	gobble.Worker
}

func NewDeliveryWorker(id int, logger *log.Logger, mailClient mail.ClientInterface, queue gobble.QueueInterface,
	globalUnsubscribesRepo models.GlobalUnsubscribesRepoInterface, unsubscribesRepo models.UnsubscribesRepoInterface,
	kindsRepo models.KindsRepoInterface, messagesRepo MessagesRepoInterface,
	database models.DatabaseInterface, sender string, encryptionKey []byte) DeliveryWorker {

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
	}

	if worker.shouldDeliver(delivery) {
		message, err := worker.pack(delivery)
		if err != nil {
			worker.logger.Printf("Not delivering because template failed to pack")
			_, err = worker.messagesRepo.Upsert(worker.database.Connection(), models.Message{
				ID:     delivery.MessageID,
				Status: StatusFailed,
			})

			if err != nil {
				worker.logger.Printf("Failed to upsert status '%s' of notification %s (failed to pack). Error: %s", StatusFailed, delivery.MessageID, err.Error())
			}

			return
		}

		status := worker.sendMail(message)
		if status != StatusDelivered {
			worker.retry(job)
			metrics.NewMetric("counter", map[string]interface{}{
				"name": "notifications.worker.retry",
			}).Log()
		} else {
			metrics.NewMetric("counter", map[string]interface{}{
				"name": "notifications.worker.delivered",
			}).Log()
		}

		_, err = worker.messagesRepo.Upsert(worker.database.Connection(), models.Message{ID: delivery.MessageID, Status: status})
		if err != nil {
			worker.logger.Printf("Failed to upsert status '%s' of notification %s to %s. Error: %s", status, delivery.MessageID, message.To, err.Error())
		}

	} else {
		metrics.NewMetric("counter", map[string]interface{}{
			"name": "notifications.worker.unsubscribed",
		}).Log()
	}

}

func (worker DeliveryWorker) retry(job *gobble.Job) {
	if job.RetryCount < 10 {
		duration := time.Duration(int64(math.Pow(2, float64(job.RetryCount))))
		job.Retry(duration * time.Minute)
		layout := "Jan 2, 2006 at 3:04pm (MST)"
		worker.logger.Printf("Message failed to send, retrying at: %s", job.ActiveAt.Format(layout))
	}
}

func (worker DeliveryWorker) shouldDeliver(delivery Delivery) bool {
	conn := worker.database.Connection()
	if worker.isCritical(conn, delivery.Options.KindID, delivery.ClientID) {
		return true
	}

	globallyUnsubscribed, err := worker.globalUnsubscribesRepo.Get(conn, delivery.UserGUID)
	if err != nil || globallyUnsubscribed {
		worker.logger.Printf("Not delivering because %s has unsubscribed", delivery.User.Emails[0])
		return false
	}

	_, err = worker.unsubscribesRepo.Find(conn, delivery.ClientID, delivery.Options.KindID, delivery.UserGUID)
	if err != nil {
		if _, ok := err.(models.RecordNotFoundError); ok {
			return len(delivery.User.Emails) > 0 && strings.Contains(delivery.User.Emails[0], "@")
		}

		worker.logger.Printf("Not delivering because: %+v", err)
		return false
	}

	worker.logger.Printf("Not delivering because %s has unsubscribed", delivery.User.Emails[0])
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

	context := NewMessageContext(delivery, worker.sender, cloak)
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
