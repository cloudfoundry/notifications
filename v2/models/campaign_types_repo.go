package models

import (
	"database/sql"
	"fmt"
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

func (r CampaignTypesRepository) Insert(connection ConnectionInterface, campaignType CampaignType) (CampaignType, error) {
	if campaignType.ID == "" {
		id, err := r.guidGenerator()
		if err != nil {
			panic(err)
		}
		campaignType.ID = id.String()
	}

	err := connection.Insert(&campaignType)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			err = DuplicateRecordError{}
		}
		return CampaignType{}, err
	}

	return campaignType, nil
}

func (r CampaignTypesRepository) GetBySenderIDAndName(connection ConnectionInterface, senderID, name string) (CampaignType, error) {
	var campaignType CampaignType
	err := connection.SelectOne(&campaignType, "SELECT * FROM `campaign_types` WHERE `sender_id` = ? AND `name` = ?", senderID, name)
	if err != nil {
		if err == sql.ErrNoRows {
			err = RecordNotFoundError{fmt.Errorf("Campaign type with sender_id %q and name %q could not be found", senderID, name)}
		}
		return CampaignType{}, err
	}

	return campaignType, nil
}

func (r CampaignTypesRepository) List(connection ConnectionInterface, senderID string) ([]CampaignType, error) {
	campaignTypeList := []CampaignType{}
	_, err := connection.Select(&campaignTypeList, "SELECT * FROM `campaign_types` WHERE `sender_id` = ?", senderID)
	return campaignTypeList, err
}

func (r CampaignTypesRepository) Get(connection ConnectionInterface, campaignTypeID string) (CampaignType, error) {
	campaignType := CampaignType{}
	err := connection.SelectOne(&campaignType, "SELECT * FROM `campaign_types` WHERE `id` = ?", campaignTypeID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = RecordNotFoundError{fmt.Errorf("Campaign type with id %q could not be found", campaignTypeID)}
		}
		return campaignType, err
	}

	return campaignType, nil
}

func (r CampaignTypesRepository) Update(connection ConnectionInterface, campaignType CampaignType) (CampaignType, error) {
	_, err := connection.Update(&campaignType)
	if err != nil {
		return CampaignType{}, err
	}

	return campaignType, err
}

func (r CampaignTypesRepository) Delete(connection ConnectionInterface, campaignType CampaignType) error {
	_, err := connection.Delete(&campaignType)
	return err
}
