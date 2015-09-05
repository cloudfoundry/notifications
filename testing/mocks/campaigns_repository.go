package mocks

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/v2/models"
)

type CampaignsRepository struct {
	InsertCall struct {
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

	ListSendingCampaignsCall struct {
		Invocations []time.Time
		Receives    struct {
			Connection models.ConnectionInterface
		}
		Returns struct {
			Campaigns []models.Campaign
			Error     error
		}
	}

	UpdateCall struct {
		Receives struct {
			Connection   models.ConnectionInterface
			CampaignList []models.Campaign
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

func (r *CampaignsRepository) Insert(conn models.ConnectionInterface, campaign models.Campaign) (models.Campaign, error) {
	r.InsertCall.Receives.Connection = conn
	r.InsertCall.Receives.Campaign = campaign

	return r.InsertCall.Returns.Campaign, r.InsertCall.Returns.Error
}

func (r *CampaignsRepository) ListSendingCampaigns(conn models.ConnectionInterface) ([]models.Campaign, error) {
	r.ListSendingCampaignsCall.Receives.Connection = conn
	r.ListSendingCampaignsCall.Invocations = append(r.ListSendingCampaignsCall.Invocations, time.Now())

	return r.ListSendingCampaignsCall.Returns.Campaigns, r.ListSendingCampaignsCall.Returns.Error
}

func (r *CampaignsRepository) Update(conn models.ConnectionInterface, campaign models.Campaign) (models.Campaign, error) {
	r.UpdateCall.Receives.Connection = conn
	r.UpdateCall.Receives.CampaignList = append(r.UpdateCall.Receives.CampaignList, campaign)

	return r.UpdateCall.Returns.Campaign, r.UpdateCall.Returns.Error
}
