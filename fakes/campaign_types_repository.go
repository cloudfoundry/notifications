package fakes

import (
	"github.com/cloudfoundry-incubator/notifications/models"
)

type CampaignTypesRepository struct {
	InsertCall struct {
		Connection         models.ConnectionInterface
		CampaignType       models.CampaignType
		ReturnCampaignType models.CampaignType
		Err                error
	}

	ListCall struct {
		Connection             models.ConnectionInterface
		ReturnCampaignTypeList []models.CampaignType
		Err                    error
	}

	GetCall struct {
		Connection     models.ConnectionInterface
		CampaignTypeID string
	}

	GetReturn struct {
		CampaignType models.CampaignType
		Err          error
	}

	UpdateCall struct {
		Connection         models.ConnectionInterface
		CampaignType       models.CampaignType
		ReturnCampaignType models.CampaignType
		Err                error
	}

	DeleteCall struct {
		Connection     models.ConnectionInterface
		CampaignTypeID string
		Err            error
	}
}

func NewCampaignTypesRepository() *CampaignTypesRepository {
	return &CampaignTypesRepository{}
}

func (r *CampaignTypesRepository) Insert(conn models.ConnectionInterface, campaignType models.CampaignType) (models.CampaignType, error) {
	r.InsertCall.CampaignType = campaignType
	r.InsertCall.Connection = conn
	return r.InsertCall.ReturnCampaignType, r.InsertCall.Err
}

func (r *CampaignTypesRepository) GetBySenderIDAndName(conn models.ConnectionInterface, senderID, name string) (models.CampaignType, error) {
	return models.CampaignType{}, nil
}

func (r *CampaignTypesRepository) List(conn models.ConnectionInterface, senderID string) ([]models.CampaignType, error) {
	r.ListCall.Connection = conn
	return r.ListCall.ReturnCampaignTypeList, r.ListCall.Err
}

func (r *CampaignTypesRepository) Get(conn models.ConnectionInterface, campaignTypeID string) (models.CampaignType, error) {
	r.GetCall.Connection = conn
	r.GetCall.CampaignTypeID = campaignTypeID
	return r.GetReturn.CampaignType, r.GetReturn.Err
}

func (r *CampaignTypesRepository) Update(conn models.ConnectionInterface, campaignType models.CampaignType) (models.CampaignType, error) {
	r.UpdateCall.Connection = conn
	r.UpdateCall.CampaignType = campaignType

	return r.UpdateCall.ReturnCampaignType, r.UpdateCall.Err
}

func (r *CampaignTypesRepository) Delete(conn models.ConnectionInterface, campaignTypeID string) error {
	r.DeleteCall.Connection = conn
	r.DeleteCall.CampaignTypeID = campaignTypeID

	return r.DeleteCall.Err
}
