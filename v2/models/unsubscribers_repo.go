package models

import (
	"database/sql"
	"fmt"
)

type Unsubscriber struct {
	ID             string `db:"id"`
	CampaignTypeID string `db:"campaign_type_id"`
	UserGUID       string `db:"user_guid"`
}

type UnsubscribersRepository struct {
	generateGUID guidGeneratorFunc
}

func NewUnsubscribersRepository(guidGenerator guidGeneratorFunc) UnsubscribersRepository {
	return UnsubscribersRepository{
		generateGUID: guidGenerator,
	}
}

func (r UnsubscribersRepository) Insert(connection ConnectionInterface, unsubscriber Unsubscriber) (Unsubscriber, error) {
	var err error
	unsubscriber.ID, err = r.generateGUID()
	if err != nil {
		return Unsubscriber{}, err
	}

	err = connection.Insert(&unsubscriber)
	if err != nil {
		return Unsubscriber{}, err
	}

	return unsubscriber, nil
}

func (r UnsubscribersRepository) Get(connection ConnectionInterface, userGUID, campaignTypeID string) (Unsubscriber, error) {
	unsubscriber := Unsubscriber{}
	err := connection.SelectOne(&unsubscriber, "SELECT * from `unsubscribers` WHERE user_guid = ? AND campaign_type_id = ?", userGUID, campaignTypeID)

	if err == sql.ErrNoRows {
		err = RecordNotFoundError{fmt.Errorf("No unsubscribers for user %s with campaign_type %s", userGUID, campaignTypeID)}
	}

	return unsubscriber, err
}

func (r UnsubscribersRepository) Delete(connection ConnectionInterface, unsubscriber Unsubscriber) error {
	_, err := connection.Exec("DELETE from `unsubscribers` WHERE user_guid = ? AND campaign_type_id = ?", unsubscriber.UserGUID, unsubscriber.CampaignTypeID)
	return err
}
