package fakes

import (
	"github.com/cloudfoundry-incubator/notifications/models"
)

type CampaignTypesRepository struct {
	InsertCall struct {
		Connection             models.ConnectionInterface
		CampaignType models.CampaignType
		ReturnCampaignType models.CampaignType
		Err                    error
	}
	ListCall struct {
		Connection                 models.ConnectionInterface
		ReturnCampaignTypeList []models.CampaignType
		Err                        error
	}
	GetCall struct {
		Connection         models.ConnectionInterface
		notificationTypeID string
	}
	GetReturn struct {
		CampaignType models.CampaignType
		Err              error
	}
}

func NewCampaignTypesRepository() *CampaignTypesRepository {
	return &CampaignTypesRepository{}
}

func (n *CampaignTypesRepository) Insert(conn models.ConnectionInterface, notificationType models.CampaignType) (models.CampaignType, error) {
	n.InsertCall.CampaignType = notificationType
	n.InsertCall.Connection = conn
	return n.InsertCall.ReturnCampaignType, n.InsertCall.Err
}

func (n *CampaignTypesRepository) GetBySenderIDAndName(conn models.ConnectionInterface, senderID, name string) (models.CampaignType, error) {
	return models.CampaignType{}, nil
}

func (n *CampaignTypesRepository) List(conn models.ConnectionInterface, senderID string) ([]models.CampaignType, error) {
	n.ListCall.Connection = conn
	return n.ListCall.ReturnCampaignTypeList, n.ListCall.Err
}

func (n *CampaignTypesRepository) Get(conn models.ConnectionInterface, notificationTypeID string) (models.CampaignType, error) {
	n.GetCall.Connection = conn
	n.GetCall.notificationTypeID = notificationTypeID
	return n.GetReturn.CampaignType, n.GetReturn.Err
}
