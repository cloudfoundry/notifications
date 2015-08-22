package collections

import "github.com/cloudfoundry-incubator/notifications/v2/models"

type campaignEnqueuer interface {
	Enqueue(campaign Campaign, jobType string) error
}

type campaignTypesGetter interface {
	Get(conn models.ConnectionInterface, campaignTypeID string) (models.CampaignType, error)
}

type templatesGetter interface {
	Get(conn models.ConnectionInterface, templateID string) (models.Template, error)
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
	templatesRepo     templatesGetter
}

func NewCampaignsCollection(enqueuer campaignEnqueuer, campaignTypesRepo campaignTypesGetter, templatesRepo templatesGetter) CampaignsCollection {
	return CampaignsCollection{
		enqueuer:          enqueuer,
		campaignTypesRepo: campaignTypesRepo,
		templatesRepo:     templatesRepo,
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

	_, err := c.templatesRepo.Get(conn, campaign.TemplateID)
	if err != nil {
		switch err.(type) {
		case models.RecordNotFoundError:
			return Campaign{}, NotFoundError{err}
		default:
			return Campaign{}, PersistenceError{err}
		}
	}

	err = c.enqueuer.Enqueue(campaign, "campaign")
	if err != nil {
		return Campaign{}, PersistenceError{Err: err}
	}

	return campaign, nil
}
