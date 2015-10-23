package v2

import (
	"strings"

	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/mail"
	"github.com/cloudfoundry-incubator/notifications/postal/common"
	"github.com/cloudfoundry-incubator/notifications/uaa"
	"github.com/cloudfoundry-incubator/notifications/v2/models"
	"github.com/pivotal-golang/lager"
)

type tokenLoader interface {
	Load(string) (string, error)
}

type messageStatusUpdater interface {
	Update(conn db.ConnectionInterface, messageID, messageStatus, campaignID string, logger lager.Logger)
}

type messagePackager interface {
	PrepareContext(delivery common.Delivery, sender, domain string) (common.MessageContext, error)
	Pack(context common.MessageContext) (mail.Message, error)
}

type userLoader interface {
	Load(userGUIDs []string, token string) (map[string]uaa.User, error)
}

type mailSender interface {
	Connect(lager.Logger) error
	Send(mail.Message, lager.Logger) error
}

type unsubscribersRepositoryInterface interface {
	Get(connection models.ConnectionInterface, userGUID, campaignTypeID string) (models.Unsubscriber, error)
}

type campaignsRepositoryInterface interface {
	Get(connection models.ConnectionInterface, campaignID string) (models.Campaign, error)
}

type metricsEmitter interface {
	Increment(counter string)
}

type DeliveryJobProcessor struct {
	mailClient              mailSender
	packager                messagePackager
	userLoader              userLoader
	tokenLoader             tokenLoader
	messageStatusUpdater    messageStatusUpdater
	unsubscribersRepository unsubscribersRepositoryInterface
	campaignsRepository     campaignsRepositoryInterface
	database                db.DatabaseInterface
	sender                  string
	domain                  string
	uaaHost                 string
	metricsEmitter          metricsEmitter
}

func NewDeliveryJobProcessor(mailClient mailSender, packager messagePackager, userLoader userLoader, tokenLoader tokenLoader,
	messageStatusUpdater messageStatusUpdater, database db.DatabaseInterface, unsubscribersRepository unsubscribersRepositoryInterface,
	campaignsRepository campaignsRepositoryInterface, sender, domain, uaaHost string, metricsEmitter metricsEmitter) DeliveryJobProcessor {

	return DeliveryJobProcessor{
		mailClient:              mailClient,
		packager:                packager,
		userLoader:              userLoader,
		tokenLoader:             tokenLoader,
		messageStatusUpdater:    messageStatusUpdater,
		campaignsRepository:     campaignsRepository,
		unsubscribersRepository: unsubscribersRepository,
		database:                database,
		sender:                  sender,
		domain:                  domain,
		uaaHost:                 uaaHost,
		metricsEmitter:          metricsEmitter,
	}
}

func (p DeliveryJobProcessor) Process(delivery common.Delivery, logger lager.Logger) error {
	conn := p.database.Connection()

	campaign, err := p.campaignsRepository.Get(conn, delivery.CampaignID)
	if err != nil {
		return err
	}

	unsubscriber, err := p.unsubscribersRepository.Get(conn, delivery.UserGUID, campaign.CampaignTypeID)
	if err != nil {
		if _, ok := err.(models.RecordNotFoundError); !ok {
			return err
		}
	}

	if unsubscriber.ID != "" {
		p.messageStatusUpdater.Update(conn, delivery.MessageID, common.StatusDelivered, delivery.CampaignID, logger)
		p.metricsEmitter.Increment("notifications.worker.unsubscribed")
		return nil
	}

	if delivery.UserGUID != "" {
		if delivery.Email != "" {
			p.messageStatusUpdater.Update(conn, delivery.MessageID, common.StatusUndeliverable, delivery.CampaignID, logger)
			return nil
		}

		token, err := p.tokenLoader.Load(p.uaaHost)
		if err != nil {
			return err
		}

		users, err := p.userLoader.Load([]string{delivery.UserGUID}, token)
		if err != nil {
			return err
		}

		emails := users[delivery.UserGUID].Emails
		if len(emails) > 0 {
			delivery.Email = emails[0]
		}
	}

	if !strings.Contains(delivery.Email, "@") {
		p.messageStatusUpdater.Update(conn, delivery.MessageID, common.StatusUndeliverable, delivery.CampaignID, logger)
		return nil
	}

	context, err := p.packager.PrepareContext(delivery, p.sender, p.domain)
	if err != nil {
		return err
	}

	message, err := p.packager.Pack(context)
	if err != nil {
		return err
	}

	err = p.mailClient.Send(message, logger)
	if err != nil {
		return err
	}

	p.messageStatusUpdater.Update(conn, delivery.MessageID, common.StatusDelivered, delivery.CampaignID, logger)

	p.metricsEmitter.Increment("notifications.worker.delivered")

	return nil
}
