package mocks

import "github.com/cloudfoundry-incubator/notifications/v2/models"

type CampaignsRepository struct {
	SetCall struct {
		Receives struct {
			Connection models.ConnectionInterface
			Campaign   models.Campaign
		}
		Returns struct {
			Campaign models.Campaign
			Error    error
		}
	}

	GetCall struct {
		Receives struct {
			Connection models.ConnectionInterface
			CampaignID string
		}
		Returns struct {
			Campaign models.Campaign
			Error    error
		}
	}
}

func NewCampaignsRepository() *CampaignsRepository {
	return &CampaignsRepository{}
}

func (r *CampaignsRepository) Get(conn models.ConnectionInterface, campaignID string) (models.Campaign, error) {
	r.GetCall.Receives.Connection = conn
	r.GetCall.Receives.CampaignID = campaignID

	return r.GetCall.Returns.Campaign, r.GetCall.Returns.Error
}

func (r *CampaignsRepository) Set(conn models.ConnectionInterface, campaign models.Campaign) (models.Campaign, error) {
	r.SetCall.Receives.Connection = conn
	r.SetCall.Receives.Campaign = campaign

	return r.SetCall.Returns.Campaign, r.SetCall.Returns.Error
}
