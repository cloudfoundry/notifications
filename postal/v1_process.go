package postal

import (
	"strings"

	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/mail"
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/v2/models"
	"github.com/pivotal-golang/lager"
)

type mailSender interface {
	Connect(lager.Logger) error
	Send(mail.Message, lager.Logger) error
}

type V1ProcessConfig struct {
	DBTrace bool
	UAAHost string
	Sender  string
	Domain  string

	Packager    Packager
	MailClient  mailSender
	Database    db.DatabaseInterface
	TokenLoader tokenLoader
	UserLoader  UserLoaderInterface

	KindsRepo              KindsRepo
	ReceiptsRepo           ReceiptsRepo
	UnsubscribesRepo       UnsubscribesRepo
	GlobalUnsubscribesRepo GlobalUnsubscribesRepo
	MessageStatusUpdater   messageStatusUpdater
	DeliveryFailureHandler deliveryFailureHandler
}

type V1Process struct {
	dbTrace bool
	uaaHost string
	sender  string
	domain  string

	packager    Packager
	mailClient  mailSender
	database    db.DatabaseInterface
	tokenLoader tokenLoader
	userLoader  UserLoaderInterface

	kindsRepo              KindsRepo
	receiptsRepo           ReceiptsRepo
	unsubscribesRepo       UnsubscribesRepo
	globalUnsubscribesRepo GlobalUnsubscribesRepo
	messageStatusUpdater   messageStatusUpdater
	deliveryFailureHandler deliveryFailureHandler
}

func NewV1Process(config V1ProcessConfig) V1Process {
	return V1Process{
		dbTrace: config.DBTrace,
		uaaHost: config.UAAHost,
		sender:  config.Sender,
		domain:  config.Domain,

		packager:    config.Packager,
		mailClient:  config.MailClient,
		database:    config.Database,
		tokenLoader: config.TokenLoader,
		userLoader:  config.UserLoader,

		kindsRepo:              config.KindsRepo,
		receiptsRepo:           config.ReceiptsRepo,
		unsubscribesRepo:       config.UnsubscribesRepo,
		globalUnsubscribesRepo: config.GlobalUnsubscribesRepo,
		messageStatusUpdater:   config.MessageStatusUpdater,
		deliveryFailureHandler: config.DeliveryFailureHandler,
	}
}

func (p V1Process) Deliver(job *gobble.Job, logger lager.Logger) error {
	var delivery Delivery
	err := job.Unmarshal(&delivery)
	if err != nil {
		metrics.NewMetric("counter", map[string]interface{}{
			"name": "notifications.worker.panic.json",
		}).Log()

		p.deliveryFailureHandler.Handle(job, logger)
		return nil
	}

	logger = logger.WithData(lager.Data{
		"message_id":      delivery.MessageID,
		"vcap_request_id": delivery.VCAPRequestID,
	})

	if p.dbTrace {
		p.database.TraceOn("", gorpCompatibleLogger{logger})
	}

	err = p.receiptsRepo.CreateReceipts(p.database.Connection(), []string{delivery.UserGUID}, delivery.ClientID, delivery.Options.KindID)
	if err != nil {
		p.deliveryFailureHandler.Handle(job, logger)
		return nil
	}

	if delivery.Email == "" {
		var token string

		token, err = p.tokenLoader.Load(p.uaaHost)
		if err != nil {
			p.deliveryFailureHandler.Handle(job, logger)
			return nil
		}

		users, err := p.userLoader.Load([]string{delivery.UserGUID}, token)
		if err != nil || len(users) < 1 {
			p.deliveryFailureHandler.Handle(job, logger)
			return nil
		}

		emails := users[delivery.UserGUID].Emails
		if len(emails) > 0 {
			delivery.Email = emails[0]
		}
	}

	logger = logger.WithData(lager.Data{
		"recipient": delivery.Email,
	})

	if p.shouldDeliver(delivery, logger) {
		status := p.deliver(delivery, logger)

		if status != StatusDelivered {
			p.deliveryFailureHandler.Handle(job, logger)
			return nil
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

	return nil
}

func (p V1Process) deliver(delivery Delivery, logger lager.Logger) string {
	context, err := p.packager.PrepareContext(delivery, p.sender, p.domain)
	if err != nil {
		panic(err)
	}

	message, err := p.packager.Pack(context)
	if err != nil {
		logger.Info("template-pack-failed")
		p.messageStatusUpdater.Update(p.database.Connection(), delivery.MessageID, StatusFailed, "", logger)
		return StatusFailed
	}

	status := p.sendMail(delivery.MessageID, message, logger)
	p.messageStatusUpdater.Update(p.database.Connection(), delivery.MessageID, status, "", logger)

	return status
}

func (p V1Process) shouldDeliver(delivery Delivery, logger lager.Logger) bool {
	conn := p.database.Connection()
	if p.isCritical(conn, delivery.Options.KindID, delivery.ClientID) {
		return true
	}

	globallyUnsubscribed, err := p.globalUnsubscribesRepo.Get(conn, delivery.UserGUID)
	if err != nil || globallyUnsubscribed {
		logger.Info("user-unsubscribed")
		p.messageStatusUpdater.Update(p.database.Connection(), delivery.MessageID, StatusUndeliverable, "", logger)
		return false
	}

	isUnsubscribed, err := p.unsubscribesRepo.Get(conn, delivery.UserGUID, delivery.ClientID, delivery.Options.KindID)
	if err != nil || isUnsubscribed {
		logger.Info("user-unsubscribed")
		p.messageStatusUpdater.Update(p.database.Connection(), delivery.MessageID, StatusUndeliverable, "", logger)
		return false
	}

	if delivery.Email == "" {
		logger.Info("no-email-address-for-user")
		p.messageStatusUpdater.Update(p.database.Connection(), delivery.MessageID, StatusUndeliverable, "", logger)
		return false
	}

	if !strings.Contains(delivery.Email, "@") {
		logger.Info("malformatted-email-address")
		p.messageStatusUpdater.Update(p.database.Connection(), delivery.MessageID, StatusUndeliverable, "", logger)
		return false
	}

	return true
}

func (p V1Process) sendMail(messageID string, message mail.Message, logger lager.Logger) string {
	err := p.mailClient.Connect(logger)
	if err != nil {
		logger.Error("smtp-connection-error", err)
		return StatusUnavailable
	}

	logger.Info("delivery-start")

	err = p.mailClient.Send(message, logger)
	if err != nil {
		logger.Error("delivery-failed-smtp-error", err)
		return StatusFailed
	}

	logger.Info("message-sent")

	return StatusDelivered
}

func (p V1Process) isCritical(conn db.ConnectionInterface, kindID, clientID string) bool {
	kind, err := p.kindsRepo.Find(conn, kindID, clientID)
	if _, ok := err.(models.RecordNotFoundError); ok {
		return false
	}

	return kind.Critical
}
