package v2

import (
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/v2/models"
	"github.com/pivotal-golang/lager"
)

type messageUpdater interface {
	Update(conn models.ConnectionInterface, message models.Message) (models.Message, error)
}

type V2MessageStatusUpdater struct {
	messages messageUpdater
}

func NewV2MessageStatusUpdater(messages messageUpdater) V2MessageStatusUpdater {
	return V2MessageStatusUpdater{
		messages: messages,
	}
}

func (mu V2MessageStatusUpdater) Update(conn db.ConnectionInterface, messageID, messageStatus, campaignID string, logger lager.Logger) {
	_, err := mu.messages.Update(conn, models.Message{
		ID:         messageID,
		Status:     messageStatus,
		CampaignID: campaignID,
	})
	if err != nil {
		logger.Session("message-updater").Error("failed-message-status-update", err, lager.Data{
			"status": messageStatus,
		})
	}
}
