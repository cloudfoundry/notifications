package v1

import (
	"strings"

	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/mail"
	"github.com/cloudfoundry-incubator/notifications/postal/common"
	"github.com/cloudfoundry-incubator/notifications/uaa"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/pivotal-golang/lager"
	"github.com/rcrowley/go-metrics"
)

type tokenLoader interface {
	Load(string) (string, error)
}

type mailSender interface {
	Connect(lager.Logger) error
	Send(mail.Message, lager.Logger) error
}

type userLoader interface {
	Load(userGUIDs []string, token string) (map[string]uaa.User, error)
}

type messageStatusUpdater interface {
	Update(conn db.ConnectionInterface, messageID, messageStatus, campaignID string, logger lager.Logger)
}

type deliveryFailureHandler interface {
	Handle(job common.Retryable, logger lager.Logger)
}

type kindsFinder interface {
	Find(connection models.ConnectionInterface, kindID string, clientID string) (models.Kind, error)
}

type receiptsCreator interface {
	CreateReceipts(connection models.ConnectionInterface, userGUIDs []string, clientID string, kindID string) error
}

type unsubscribesGetter interface {
	Get(connection models.ConnectionInterface, userGUID string, clientID string, kindID string) (bool, error)
}

type globalUnsubscribesGetter interface {
	Get(connection models.ConnectionInterface, userGUID string) (bool, error)
}

type DeliveryJobProcessorConfig struct {
	DBTrace bool
	UAAHost string
	Sender  string
	Domain  string

	Packager    common.Packager
	MailClient  mailSender
	Database    db.DatabaseInterface
	TokenLoader tokenLoader
	UserLoader  userLoader

	KindsRepo              kindsFinder
	ReceiptsRepo           receiptsCreator
	UnsubscribesRepo       unsubscribesGetter
	GlobalUnsubscribesRepo globalUnsubscribesGetter
	MessageStatusUpdater   messageStatusUpdater
	DeliveryFailureHandler deliveryFailureHandler
}

type DeliveryJobProcessor struct {
	dbTrace bool
	uaaHost string
	sender  string
	domain  string

	packager    common.Packager
	mailClient  mailSender
	database    db.DatabaseInterface
	tokenLoader tokenLoader
	userLoader  userLoader

	kindsRepo              kindsFinder
	receiptsRepo           receiptsCreator
	unsubscribesRepo       unsubscribesGetter
	globalUnsubscribesRepo globalUnsubscribesGetter
	messageStatusUpdater   messageStatusUpdater
	deliveryFailureHandler deliveryFailureHandler
}

func NewDeliveryJobProcessor(config DeliveryJobProcessorConfig) DeliveryJobProcessor {
	return DeliveryJobProcessor{
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

func (p DeliveryJobProcessor) Process(job *gobble.Job, logger lager.Logger) error {
	var delivery common.Delivery
	err := job.Unmarshal(&delivery)
	if err != nil {
		metrics.GetOrRegisterCounter("notifications.worker.panic.json", nil).Inc(1)

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
		status := p.process(delivery, logger)

		if status != common.StatusDelivered {
			p.deliveryFailureHandler.Handle(job, logger)
			return nil
		} else {
			metrics.GetOrRegisterCounter("notifications.worker.delivered", nil).Inc(1)
		}
	} else {
		metrics.GetOrRegisterCounter("notifications.worker.unsubscribed", nil).Inc(1)
	}

	return nil
}

func (p DeliveryJobProcessor) process(delivery common.Delivery, logger lager.Logger) string {
	context, err := p.packager.PrepareContext(delivery, p.sender, p.domain)
	if err != nil {
		panic(err)
	}

	message, err := p.packager.Pack(context)
	if err != nil {
		logger.Info("template-pack-failed")
		p.messageStatusUpdater.Update(p.database.Connection(), delivery.MessageID, common.StatusFailed, "", logger)
		return common.StatusFailed
	}

	status := p.sendMail(delivery.MessageID, message, logger)
	p.messageStatusUpdater.Update(p.database.Connection(), delivery.MessageID, status, "", logger)

	return status
}

func (p DeliveryJobProcessor) shouldDeliver(delivery common.Delivery, logger lager.Logger) bool {
	conn := p.database.Connection()
	if p.isCritical(conn, delivery.Options.KindID, delivery.ClientID) {
		return true
	}

	globallyUnsubscribed, err := p.globalUnsubscribesRepo.Get(conn, delivery.UserGUID)
	if err != nil || globallyUnsubscribed {
		logger.Info("user-unsubscribed")
		p.messageStatusUpdater.Update(p.database.Connection(), delivery.MessageID, common.StatusUndeliverable, "", logger)
		return false
	}

	isUnsubscribed, err := p.unsubscribesRepo.Get(conn, delivery.UserGUID, delivery.ClientID, delivery.Options.KindID)
	if err != nil || isUnsubscribed {
		logger.Info("user-unsubscribed")
		p.messageStatusUpdater.Update(p.database.Connection(), delivery.MessageID, common.StatusUndeliverable, "", logger)
		return false
	}

	if delivery.Email == "" {
		logger.Info("no-email-address-for-user")
		p.messageStatusUpdater.Update(p.database.Connection(), delivery.MessageID, common.StatusUndeliverable, "", logger)
		return false
	}

	if !strings.Contains(delivery.Email, "@") {
		logger.Info("malformatted-email-address")
		p.messageStatusUpdater.Update(p.database.Connection(), delivery.MessageID, common.StatusUndeliverable, "", logger)
		return false
	}

	return true
}

func (p DeliveryJobProcessor) sendMail(messageID string, message mail.Message, logger lager.Logger) string {
	err := p.mailClient.Connect(logger)
	if err != nil {
		logger.Error("smtp-connection-error", err)
		return common.StatusFailed
	}

	logger.Info("delivery-start")

	err = p.mailClient.Send(message, logger)
	if err != nil {
		logger.Error("delivery-failed-smtp-error", err)
		return common.StatusFailed
	}

	logger.Info("message-sent")

	return common.StatusDelivered
}

func (p DeliveryJobProcessor) isCritical(conn db.ConnectionInterface, kindID, clientID string) bool {
	kind, err := p.kindsRepo.Find(conn, kindID, clientID)
	if _, ok := err.(models.NotFoundError); ok {
		return false
	}

	return kind.Critical
}
