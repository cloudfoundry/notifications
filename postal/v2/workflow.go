package v2

import (
	"strings"

	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/mail"
	"github.com/cloudfoundry-incubator/notifications/postal/common"
	"github.com/cloudfoundry-incubator/notifications/uaa"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
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

type Workflow struct {
	mailClient              mailSender
	packager                messagePackager
	userLoader              userLoader
	tokenLoader             tokenLoader
	messageStatusUpdater    messageStatusUpdater
	unsubscribersRepository unsubscribersRepositoryInterface
	campaignsRepository     campaignsRepositoryInterface
	database                services.DatabaseInterface
	sender                  string
	domain                  string
	uaaHost                 string
}

func NewWorkflow(mailClient mailSender, packager messagePackager, userLoader userLoader, tokenLoader tokenLoader,
	messageStatusUpdater messageStatusUpdater, database services.DatabaseInterface, unsubscribersRepository unsubscribersRepositoryInterface,
	campaignsRepository campaignsRepositoryInterface, sender, domain, uaaHost string) Workflow {
	return Workflow{
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
	}
}

func (w Workflow) Deliver(delivery common.Delivery, logger lager.Logger) error {
	conn := w.database.Connection()

	campaign, err := w.campaignsRepository.Get(conn, delivery.CampaignID)
	if err != nil {
		return err
	}

	unsubscriber, err := w.unsubscribersRepository.Get(conn, delivery.UserGUID, campaign.CampaignTypeID)
	if err != nil {
		if _, ok := err.(models.RecordNotFoundError); !ok {
			return err
		}
	}

	if unsubscriber.ID != "" {
		w.messageStatusUpdater.Update(w.database.Connection(), delivery.MessageID, common.StatusDelivered, delivery.CampaignID, logger)
		return nil
	}

	token, err := w.tokenLoader.Load(w.uaaHost)
	if err != nil {
		return err
	}

	users, err := w.userLoader.Load([]string{delivery.UserGUID}, token)
	if err != nil {
		return err
	}

	emails := users[delivery.UserGUID].Emails
	if len(emails) > 0 {
		delivery.Email = emails[0]
	}

	if !strings.Contains(delivery.Email, "@") {
		w.messageStatusUpdater.Update(w.database.Connection(), delivery.MessageID, common.StatusUndeliverable, delivery.CampaignID, logger)
		return nil
	}

	context, err := w.packager.PrepareContext(delivery, w.sender, w.domain)
	if err != nil {
		return err
	}

	message, err := w.packager.Pack(context)
	if err != nil {
		return err
	}

	err = w.mailClient.Send(message, logger)
	if err != nil {
		return err
	}

	w.messageStatusUpdater.Update(w.database.Connection(), delivery.MessageID, common.StatusDelivered, delivery.CampaignID, logger)

	return nil
}
