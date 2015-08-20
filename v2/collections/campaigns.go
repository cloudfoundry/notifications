package collections

type campaignEnqueuer interface {
	Enqueue(campaign Campaign, jobType string) error
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
	enqueuer campaignEnqueuer
}

func NewCampaignsCollection(enqueuer campaignEnqueuer) CampaignsCollection {
	return CampaignsCollection{
		enqueuer: enqueuer,
	}
}

func (c CampaignsCollection) Create(conn ConnectionInterface, campaign Campaign) (Campaign, error) {
	campaign.ID = "some-random-id"

	err := c.enqueuer.Enqueue(campaign, "campaign")
	if err != nil {
		return campaign, PersistenceError{Err: err}
	}

	return campaign, nil
}
