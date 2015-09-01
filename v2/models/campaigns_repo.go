package models

type CampaignsRepository struct {
	guidGenerator guidGeneratorFunc
}

func NewCampaignsRepository(guidGenerator guidGeneratorFunc) CampaignsRepository {
	return CampaignsRepository{
		guidGenerator: guidGenerator,
	}
}

func (r CampaignsRepository) Set(conn ConnectionInterface, campaign Campaign) (Campaign, error) {
	id, err := r.guidGenerator()
	if err != nil {
		panic(err)
	}

	campaign.ID = id.String()

	err = conn.Insert(&campaign)
	if err != nil {
		panic(err)
	}

	return campaign, nil
}

func (r CampaignsRepository) Get(conn ConnectionInterface, campaignID string) (Campaign, error) {
	campaign := Campaign{}

	err := conn.SelectOne(&campaign, "SELECT * FROM `campaigns` WHERE `id` = ?", campaignID)
	if err != nil {
		panic(err)
	}

	return campaign, nil
}
