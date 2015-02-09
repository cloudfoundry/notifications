package services

import (
	"fmt"

	"github.com/cloudfoundry-incubator/notifications/models"
)

type PreferenceUpdaterInterface interface {
	Execute(models.ConnectionInterface, []models.Preference, bool, string) error
}

type PreferenceUpdater struct {
	globalUnsubscribesRepo models.GlobalUnsubscribesRepoInterface
	unsubscribesRepo       models.UnsubscribesRepoInterface
	kindsRepo              models.KindsRepoInterface
}

func NewPreferenceUpdater(globalUnsubscribesRepo models.GlobalUnsubscribesRepoInterface, unsubscribesRepo models.UnsubscribesRepoInterface, kindsRepo models.KindsRepoInterface) PreferenceUpdater {
	return PreferenceUpdater{
		globalUnsubscribesRepo: globalUnsubscribesRepo,
		unsubscribesRepo:       unsubscribesRepo,
		kindsRepo:              kindsRepo,
	}
}

func (updater PreferenceUpdater) Execute(conn models.ConnectionInterface, preferences []models.Preference, globalUnsubscribe bool, userID string) error {
	err := updater.globalUnsubscribesRepo.Set(conn, userID, globalUnsubscribe)
	if err != nil {
		return err
	}

	for _, preference := range preferences {

		kind, err := updater.kindsRepo.Find(conn, preference.KindID, preference.ClientID)
		if err != nil {
			return MissingKindOrClientError(fmt.Sprintf("The kind '%s' cannot be found for client '%s'", preference.KindID, preference.ClientID))
		}

		if kind.Critical {
			return CriticalKindError(fmt.Sprintf("The kind '%s' for the '%s' client is critical and cannot be unsubscribed from", preference.KindID, preference.ClientID))
		}

		if !preference.Email {
			_, err := updater.unsubscribesRepo.Upsert(conn, models.Unsubscribe{
				ClientID: preference.ClientID,
				KindID:   preference.KindID,
				UserID:   userID,
			})
			if err != nil {
				return err
			}
		} else {
			_, err := updater.unsubscribesRepo.Destroy(conn, models.Unsubscribe{
				ClientID: preference.ClientID,
				KindID:   preference.KindID,
				UserID:   userID,
			})
			if err != nil {
				return err
			}

		}
	}
	return nil
}
