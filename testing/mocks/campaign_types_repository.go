package mocks

import "github.com/cloudfoundry-incubator/notifications/v2/models"

type CampaignTypesRepository struct {
	InsertCall struct {
		Receives struct {
			Connection   models.ConnectionInterface
			CampaignType models.CampaignType
		}
		Returns struct {
			CampaignType models.CampaignType
			Error        error
		}
	}

	ListCall struct {
		Receives struct {
			Connection models.ConnectionInterface
			SenderID   string
		}
		Returns struct {
			CampaignTypeList []models.CampaignType
			Error            error
		}
	}

	GetCall struct {
		Receives struct {
			Connection     models.ConnectionInterface
			CampaignTypeID string
		}
		Returns struct {
			CampaignType models.CampaignType
			Error        error
		}
	}

	UpdateCall struct {
		Receives struct {
			Connection   models.ConnectionInterface
			CampaignType models.CampaignType
		}
		Returns struct {
			CampaignType models.CampaignType
			Error        error
		}
	}

	DeleteCall struct {
		Receives struct {
			Connection   models.ConnectionInterface
			CampaignType models.CampaignType
		}
		Returns struct {
			Error error
		}
	}
}

func NewCampaignTypesRepository() *CampaignTypesRepository {
	return &CampaignTypesRepository{}
}

func (r *CampaignTypesRepository) Insert(conn models.ConnectionInterface, campaignType models.CampaignType) (models.CampaignType, error) {
	r.InsertCall.Receives.CampaignType = campaignType
	r.InsertCall.Receives.Connection = conn

	return r.InsertCall.Returns.CampaignType, r.InsertCall.Returns.Error
}

func (r *CampaignTypesRepository) GetBySenderIDAndName(conn models.ConnectionInterface, senderID, name string) (models.CampaignType, error) {
	return models.CampaignType{}, nil
}

func (r *CampaignTypesRepository) List(conn models.ConnectionInterface, senderID string) ([]models.CampaignType, error) {
	r.ListCall.Receives.Connection = conn
	r.ListCall.Receives.SenderID = senderID

	return r.ListCall.Returns.CampaignTypeList, r.ListCall.Returns.Error
}

func (r *CampaignTypesRepository) Get(conn models.ConnectionInterface, campaignTypeID string) (models.CampaignType, error) {
	r.GetCall.Receives.Connection = conn
	r.GetCall.Receives.CampaignTypeID = campaignTypeID

	return r.GetCall.Returns.CampaignType, r.GetCall.Returns.Error
}

func (r *CampaignTypesRepository) Update(conn models.ConnectionInterface, campaignType models.CampaignType) (models.CampaignType, error) {
	r.UpdateCall.Receives.Connection = conn
	r.UpdateCall.Receives.CampaignType = campaignType

	return r.UpdateCall.Returns.CampaignType, r.UpdateCall.Returns.Error
}

func (r *CampaignTypesRepository) Delete(conn models.ConnectionInterface, campaignType models.CampaignType) error {
	r.DeleteCall.Receives.Connection = conn
	r.DeleteCall.Receives.CampaignType = campaignType

	return r.DeleteCall.Returns.Error
}
