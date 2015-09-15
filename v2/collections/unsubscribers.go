package collections

import (
	"fmt"

	"github.com/cloudfoundry-incubator/notifications/v2/models"
)

type Unsubscriber struct {
	ID             string
	CampaignTypeID string
	UserGUID       string
}

type unsubscribersSetter interface {
	Insert(connection models.ConnectionInterface, unsubscriber models.Unsubscriber) (models.Unsubscriber, error)
}

type UnsubscribersCollection struct {
	unsubscribersRepository unsubscribersSetter
	userFinder              existenceChecker
	campaignTypesRepository campaignTypesGetter
}

func NewUnsubscribersCollection(unsubscribersRepository unsubscribersSetter,
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
		return Unsubscriber{}, PersistenceError{err}
	}

	unsubscriber.ID = unsub.ID
	return unsubscriber, nil
}
