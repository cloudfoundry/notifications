package mocks

import (
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/pivotal-golang/lager"
)

type MessageStatusUpdater struct {
	UpdateCall struct {
		Receives struct {
			Connection    db.ConnectionInterface
			MessageID     string
			MessageStatus string
			CampaignID    string
			Logger        lager.Logger
		}
	}
}

func NewMessageStatusUpdater() *MessageStatusUpdater {
	return &MessageStatusUpdater{}
}

func (msu *MessageStatusUpdater) Update(conn db.ConnectionInterface, messageID, messageStatus, campaignID string, logger lager.Logger) {
	msu.UpdateCall.Receives.Connection = conn
	msu.UpdateCall.Receives.MessageID = messageID
	msu.UpdateCall.Receives.MessageStatus = messageStatus
	msu.UpdateCall.Receives.CampaignID = campaignID
	msu.UpdateCall.Receives.Logger = logger
}
