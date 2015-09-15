package mocks

import (
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/v2/models"
)

type UnsubscribersRepository struct {
	InsertCall struct {
		Receives struct {
			Unsubscriber models.Unsubscriber
			Connection   db.ConnectionInterface
		}
		Returns struct {
			Unsubscriber models.Unsubscriber
			Error        error
		}
	}
	GetCall struct {
		Receives struct {
			CampaignTypeID string
			UserGUID       string
			Connection     db.ConnectionInterface
		}
		Returns struct {
			Unsubscriber models.Unsubscriber
			Error        error
		}
	}
}

func NewUnsubscribersRepository() *UnsubscribersRepository {
	return &UnsubscribersRepository{}
}

func (ur *UnsubscribersRepository) Insert(connection models.ConnectionInterface, unsubscriber models.Unsubscriber) (models.Unsubscriber, error) {
	ur.InsertCall.Receives.Connection = connection
	ur.InsertCall.Receives.Unsubscriber = unsubscriber
	return ur.InsertCall.Returns.Unsubscriber, ur.InsertCall.Returns.Error
}

func (ur *UnsubscribersRepository) Get(connection models.ConnectionInterface, userGUID, campaignTypeID string) (models.Unsubscriber, error) {
	ur.GetCall.Receives.CampaignTypeID = campaignTypeID
	ur.GetCall.Receives.UserGUID = userGUID
	ur.GetCall.Receives.Connection = connection

	return ur.GetCall.Returns.Unsubscriber, ur.GetCall.Returns.Error
}
