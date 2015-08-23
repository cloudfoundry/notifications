package mocks

import (
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/pivotal-golang/lager"
)

type MessageStatusUpdater struct {
	UpdateCall struct {
		Receives struct {
			Connection    models.ConnectionInterface
			MessageID     string
			MessageStatus string
			Logger        lager.Logger
		}
	}
}

func NewMessageStatusUpdater() *MessageStatusUpdater {
	return &MessageStatusUpdater{}
}

func (msu *MessageStatusUpdater) Update(conn models.ConnectionInterface, messageID, messageStatus string, logger lager.Logger) {
	msu.UpdateCall.Receives.Connection = conn
	msu.UpdateCall.Receives.MessageID = messageID
	msu.UpdateCall.Receives.MessageStatus = messageStatus
	msu.UpdateCall.Receives.Logger = logger
}
