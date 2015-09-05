package models

import (
	"database/sql"
	"fmt"
)

type CampaignsRepository struct {
	guidGenerator guidGeneratorFunc
}

func NewCampaignsRepository(guidGenerator guidGeneratorFunc) CampaignsRepository {
	return CampaignsRepository{
		guidGenerator: guidGenerator,
	}
}

func (r CampaignsRepository) Insert(conn ConnectionInterface, campaign Campaign) (Campaign, error) {
	id, err := r.guidGenerator()
	if err != nil {
		panic(err)
	}

	campaign.ID = id.String()

	err = conn.Insert(&campaign)
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (r CampaignsRepository) Get(conn ConnectionInterface, campaignID string) (Campaign, error) {
	campaign := Campaign{}

	err := conn.SelectOne(&campaign, "SELECT * FROM `campaigns` WHERE `id` = ?", campaignID)
	if err != nil {
		if err == sql.ErrNoRows {
			return campaign, RecordNotFoundError{fmt.Errorf("Campaign with id %q could not be found", campaignID)}
		}

		return campaign, err
	}

	return campaign, nil
}

func (r CampaignsRepository) ListSendingCampaigns(conn ConnectionInterface) ([]Campaign, error) {
	campaignList := []Campaign{}

	_, err := conn.Select(&campaignList, "SELECT * FROM `campaigns` WHERE `status` != \"completed\"")

	return campaignList, err
}
