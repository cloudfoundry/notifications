package models

import (
	"database/sql"
	"strings"
)

type CampaignTypesRepository struct {
	guidGenerator guidGeneratorFunc
}

func NewCampaignTypesRepository(guidGenerator guidGeneratorFunc) CampaignTypesRepository {
	return CampaignTypesRepository{
		guidGenerator: guidGenerator,
	}
}

func (n CampaignTypesRepository) Insert(connection ConnectionInterface, campaignType CampaignType) (CampaignType, error) {
	id, err := n.guidGenerator()
	if err != nil {
		panic(err)
	}

	campaignType.ID = id.String()
	err = connection.Insert(&campaignType)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			err = DuplicateRecordError{}
		}
		return campaignType, err
	}

	return campaignType, nil
}

func (n CampaignTypesRepository) GetBySenderIDAndName(connection ConnectionInterface, senderID, name string) (CampaignType, error) {
	campaignType := CampaignType{}
	err := connection.SelectOne(&campaignType, "SELECT * FROM `campaign_types` WHERE `sender_id` = ? AND `name` = ?", senderID, name)
	if err != nil {
		if err == sql.ErrNoRows {
			err = NewRecordNotFoundError("Campaign type with sender_id %q and name %q could not be found", senderID, name)
		}
		return campaignType, err
	}

	return campaignType, nil
}

func (n CampaignTypesRepository) List(connection ConnectionInterface, senderID string) ([]CampaignType, error) {
	campaignTypeList := []CampaignType{}
	_, err := connection.Select(&campaignTypeList, "SELECT * FROM `campaign_types` WHERE `sender_id` = ?", senderID)
	return campaignTypeList, err
}

func (n CampaignTypesRepository) Get(connection ConnectionInterface, campaignTypeID string) (CampaignType, error) {
	campaignType := CampaignType{}
	err := connection.SelectOne(&campaignType, "SELECT * FROM `campaign_types` WHERE `id` = ?", campaignTypeID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = NewRecordNotFoundError("Campaign type with id %q could not be found", campaignTypeID)
		}
		return campaignType, err
	}

	return campaignType, nil
}
