package handlers

import "github.com/cloudfoundry-incubator/notifications/models"

type Preference struct {
    repo models.PreferencesRepoInterface
}

type NotificationPreferences map[string]map[string]map[string]string

type PreferenceInterface interface {
    Execute(string) (NotificationPreferences, error)
}

func NewPreference(repo models.PreferencesRepoInterface) *Preference {
    return &Preference{
        repo: repo,
    }
}

func (preference Preference) Execute(UserGUID string) (NotificationPreferences, error) {
    data, err := preference.repo.FindNonCriticalPreferences(models.Database().Connection, UserGUID)
    if err != nil {
        return NotificationPreferences{}, err
    }

    resultsMap := NotificationPreferences{}
    for _, preferenceData := range data {

        if userPreferences, ok := resultsMap[preferenceData.ClientID]; ok {
            userPreferences[preferenceData.KindID] = map[string]string{
                "email": "true",
            }
        } else {
            resultsMap[preferenceData.ClientID] = map[string]map[string]string{
                preferenceData.KindID: map[string]string{
                    "email": "true",
                },
            }
        }

    }
    return resultsMap, nil
}
