package collections

import (
	"fmt"

	"github.com/cloudfoundry-incubator/notifications/v2/models"
)

type existenceChecker interface {
	Exists(guid string) (bool, error)
}

type Unsubscriber struct {
	ID             string
	CampaignTypeID string
	UserGUID       string
}

type unsubscribersSetterDeleter interface {
	Insert(connection models.ConnectionInterface, unsubscriber models.Unsubscriber) (models.Unsubscriber, error)
	Delete(connection models.ConnectionInterface, unsubscriber models.Unsubscriber) error
}

type UnsubscribersCollection struct {
	unsubscribersRepository unsubscribersSetterDeleter
	userFinder              existenceChecker
	campaignTypesRepository campaignTypesGetter
}

func NewUnsubscribersCollection(unsubscribersRepository unsubscribersSetterDeleter,
	campaignTypesRepository campaignTypesGetter, userFinder existenceChecker) UnsubscribersCollection {

	return UnsubscribersCollection{
		unsubscribersRepository: unsubscribersRepository,
		userFinder:              userFinder,
		campaignTypesRepository: campaignTypesRepository,
	}
}

func (c UnsubscribersCollection) Set(connection ConnectionInterface, unsubscriber Unsubscriber) (Unsubscriber, error) {
	campaignType, err := c.campaignTypesRepository.Get(connection, unsubscriber.CampaignTypeID)
	if err != nil {
		if e, ok := err.(models.RecordNotFoundError); ok {
			return Unsubscriber{}, NotFoundError{e}
		}
		return Unsubscriber{}, err
	}

	if campaignType.Critical {
		return Unsubscriber{}, PermissionsError{fmt.Errorf("Campaign type %q cannot be unsubscribed from", unsubscriber.CampaignTypeID)}
	}

	userExists, err := c.userFinder.Exists(unsubscriber.UserGUID)
	if err != nil {
		return Unsubscriber{}, err
	}

	if !userExists {
		return Unsubscriber{}, NotFoundError{fmt.Errorf("User %q not found", unsubscriber.UserGUID)}
	}

	unsub, err := c.unsubscribersRepository.Insert(connection, models.Unsubscriber{
		CampaignTypeID: unsubscriber.CampaignTypeID,
		UserGUID:       unsubscriber.UserGUID,
	})
	if err != nil {
		return Unsubscriber{}, err
	}

	unsubscriber.ID = unsub.ID
	return unsubscriber, nil
}

func (c UnsubscribersCollection) Delete(connection ConnectionInterface, unsubscriber Unsubscriber) error {
	_, err := c.campaignTypesRepository.Get(connection, unsubscriber.CampaignTypeID)
	if err != nil {
		if e, ok := err.(models.RecordNotFoundError); ok {
			return NotFoundError{e}
		}
		return err
	}

	userExists, err := c.userFinder.Exists(unsubscriber.UserGUID)
	if err != nil {
		return err
	}

	if !userExists {
		return NotFoundError{fmt.Errorf("User %q not found", unsubscriber.UserGUID)}
	}

	err = c.unsubscribersRepository.Delete(connection, models.Unsubscriber{
		CampaignTypeID: unsubscriber.CampaignTypeID,
		UserGUID:       unsubscriber.UserGUID,
	})

	return err
}
