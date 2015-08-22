package collections

import "github.com/cloudfoundry-incubator/notifications/v2/models"

type campaignEnqueuer interface {
	Enqueue(campaign Campaign, jobType string) error
}

type campaignTypesGetter interface {
	Get(conn models.ConnectionInterface, campaignTypeID string) (models.CampaignType, error)
}

type Campaign struct {
	ID             string
	SendTo         map[string]string
	CampaignTypeID string
	Text           string
	HTML           string
	Subject        string
	TemplateID     string
	ReplyTo        string
	ClientID       string
}

type CampaignsCollection struct {
	enqueuer          campaignEnqueuer
	campaignTypesRepo campaignTypesGetter
}

func NewCampaignsCollection(enqueuer campaignEnqueuer, campaignTypesRepo campaignTypesGetter) CampaignsCollection {
	return CampaignsCollection{
		enqueuer:          enqueuer,
		campaignTypesRepo: campaignTypesRepo,
	}
}

func (c CampaignsCollection) Create(conn ConnectionInterface, campaign Campaign) (Campaign, error) {
	campaign.ID = "some-random-id"

	if campaign.TemplateID == "" {
		campaignType, err := c.campaignTypesRepo.Get(conn, campaign.CampaignTypeID)
		if err != nil {
			panic(err)
		}

		campaign.TemplateID = campaignType.TemplateID
	}

	err := c.enqueuer.Enqueue(campaign, "campaign")
	if err != nil {
		return campaign, PersistenceError{Err: err}
	}

	return campaign, nil
}
