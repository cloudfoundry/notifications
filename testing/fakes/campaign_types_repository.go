package fakes

import (
	"github.com/cloudfoundry-incubator/notifications/models"
)

type CampaignTypesRepository struct {
	InsertCall struct {
		Receives struct {
			Connection   models.ConnectionInterface
			CampaignType models.CampaignType
		}
		Returns struct {
			CampaignType models.CampaignType
			Err          error
		}
	}

	ListCall struct {
		Receives struct {
			Connection models.ConnectionInterface
		}
		Returns struct {
			CampaignTypeList []models.CampaignType
			Err              error
		}
	}

	GetCall struct {
		Receives struct {
			Connection     models.ConnectionInterface
			CampaignTypeID string
		}
		Returns struct {
			CampaignType models.CampaignType
			Err          error
		}
	}

	UpdateCall struct {
		Receives struct {
			Connection   models.ConnectionInterface
			CampaignType models.CampaignType
		}
		Returns struct {
			CampaignType models.CampaignType
			Err          error
		}
	}

	DeleteCall struct {
		Receives struct {
			Connection   models.ConnectionInterface
			CampaignType models.CampaignType
		}
		Returns struct {
			Err error
		}
	}
}

func NewCampaignTypesRepository() *CampaignTypesRepository {
	return &CampaignTypesRepository{}
}

func (r *CampaignTypesRepository) Insert(conn models.ConnectionInterface, campaignType models.CampaignType) (models.CampaignType, error) {
	r.InsertCall.Receives.CampaignType = campaignType
	r.InsertCall.Receives.Connection = conn

	return r.InsertCall.Returns.CampaignType, r.InsertCall.Returns.Err
}

func (r *CampaignTypesRepository) GetBySenderIDAndName(conn models.ConnectionInterface, senderID, name string) (models.CampaignType, error) {
	return models.CampaignType{}, nil
}

func (r *CampaignTypesRepository) List(conn models.ConnectionInterface, senderID string) ([]models.CampaignType, error) {
	r.ListCall.Receives.Connection = conn

	return r.ListCall.Returns.CampaignTypeList, r.ListCall.Returns.Err
}

func (r *CampaignTypesRepository) Get(conn models.ConnectionInterface, campaignTypeID string) (models.CampaignType, error) {
	r.GetCall.Receives.Connection = conn
	r.GetCall.Receives.CampaignTypeID = campaignTypeID

	return r.GetCall.Returns.CampaignType, r.GetCall.Returns.Err
}

func (r *CampaignTypesRepository) Update(conn models.ConnectionInterface, campaignType models.CampaignType) (models.CampaignType, error) {
	r.UpdateCall.Receives.Connection = conn
	r.UpdateCall.Receives.CampaignType = campaignType

	return r.UpdateCall.Returns.CampaignType, r.UpdateCall.Returns.Err
}

func (r *CampaignTypesRepository) Delete(conn models.ConnectionInterface, campaignType models.CampaignType) error {
	r.DeleteCall.Receives.Connection = conn
	r.DeleteCall.Receives.CampaignType = campaignType

	return r.DeleteCall.Returns.Err
}
