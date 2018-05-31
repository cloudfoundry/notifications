package v1

import (
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/pivotal-golang/lager"
)

type MessageStatusUpdater struct {
	messagesRepo MessageUpserter
}

type MessageUpserter interface {
	Upsert(conn models.ConnectionInterface, message models.Message) (models.Message, error)
}

func NewMessageStatusUpdater(messagesRepo MessageUpserter) MessageStatusUpdater {
	return MessageStatusUpdater{
		messagesRepo: messagesRepo,
	}
}

func (mu MessageStatusUpdater) Update(conn db.ConnectionInterface, messageID, messageStatus, campaignID string, logger lager.Logger) {
	_, err := mu.messagesRepo.Upsert(conn, models.Message{
		ID:         messageID,
		Status:     messageStatus,
	})
	if err != nil {
		logger.Session("message-updater").Error("failed-message-status-upsert", err, lager.Data{
			"status": messageStatus,
		})
	}
}
