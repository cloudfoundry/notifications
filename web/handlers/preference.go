package handlers

import "github.com/cloudfoundry-incubator/notifications/models"

type Preference struct {
    repo models.PreferencesRepoInterface
}

type PreferenceInterface interface {
    Execute(string) (NotificationPreferences, error)
}

func NewPreference(repo models.PreferencesRepoInterface) *Preference {
    return &Preference{
        repo: repo,
    }
}

func (preference Preference) Execute(UserGUID string) (NotificationPreferences, error) {
    preferencesData, err := preference.repo.FindNonCriticalPreferences(models.Database().Connection, UserGUID)

    if err != nil {
        return NotificationPreferences{}, err
    }

    preferences := NewNotificationPreferences()

    for _, preferenceData := range preferencesData {

        //true will change after story #76994664
        preferences.Add(preferenceData.ClientID, preferenceData.KindID, true)

    }
    return preferences, nil
}
