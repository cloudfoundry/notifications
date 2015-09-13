package services

import (
	"fmt"

	"github.com/cloudfoundry-incubator/notifications/v1/models"
)

type PreferenceUpdater struct {
	globalUnsubscribesRepo GlobalUnsubscribesRepo
	unsubscribesRepo       UnsubscribesRepo
	kindsRepo              KindsRepo
}

func NewPreferenceUpdater(globalUnsubscribesRepo GlobalUnsubscribesRepo, unsubscribesRepo UnsubscribesRepo, kindsRepo KindsRepo) PreferenceUpdater {
	return PreferenceUpdater{
		globalUnsubscribesRepo: globalUnsubscribesRepo,
		unsubscribesRepo:       unsubscribesRepo,
		kindsRepo:              kindsRepo,
	}
}

func (updater PreferenceUpdater) Update(conn ConnectionInterface, preferences []models.Preference, globalUnsubscribe bool, userID string) error {
	err := updater.globalUnsubscribesRepo.Set(conn, userID, globalUnsubscribe)
	if err != nil {
		return err
	}

	for _, preference := range preferences {
		kind, err := updater.kindsRepo.Find(conn, preference.KindID, preference.ClientID)
		if err != nil {
			return MissingKindOrClientError{fmt.Errorf("The kind '%s' cannot be found for client '%s'", preference.KindID, preference.ClientID)}
		}

		if kind.Critical {
			return CriticalKindError{fmt.Errorf("The kind '%s' for the '%s' client is critical and cannot be unsubscribed from", preference.KindID, preference.ClientID)}
		}

		err = updater.unsubscribesRepo.Set(conn, userID, preference.ClientID, preference.KindID, !preference.Email)
		if err != nil {
			return err
		}
	}
	return nil
}
