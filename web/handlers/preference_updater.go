package handlers

import "github.com/cloudfoundry-incubator/notifications/models"

type PreferenceUpdaterInterface interface {
    Execute(models.ConnectionInterface, []models.Preference, string) error
}

type PreferenceUpdater struct {
    repo models.UnsubscribesRepoInterface
}

func NewPreferenceUpdater(repo models.UnsubscribesRepoInterface) PreferenceUpdater {
    return PreferenceUpdater{
        repo: repo,
    }
}

func (updater PreferenceUpdater) Execute(conn models.ConnectionInterface, preferences []models.Preference, userID string) error {
    for _, preference := range preferences {
        _, err := updater.repo.Upsert(conn, models.Unsubscribe{
            ClientID: preference.ClientID,
            KindID:   preference.KindID,
            UserID:   userID,
        })
        //TODO: error handling
        if err != nil {
            return err
        }
    }
    return nil
}
