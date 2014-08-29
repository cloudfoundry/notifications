package services

import "github.com/cloudfoundry-incubator/notifications/models"

type Preference struct {
    repo models.PreferencesRepoInterface
}

type PreferenceInterface interface {
    Execute(string) (PreferencesBuilder, error)
}

func NewPreference(repo models.PreferencesRepoInterface) *Preference {
    return &Preference{
        repo: repo,
    }
}

func (preference Preference) Execute(UserGUID string) (PreferencesBuilder, error) {
    preferencesData, err := preference.repo.FindNonCriticalPreferences(models.Database().Connection, UserGUID)
    if err != nil {
        return PreferencesBuilder{}, err
    }

    preferences := NewPreferencesBuilder()
    for _, preferenceData := range preferencesData {
        preferences.Add(preferenceData.ClientID, preferenceData.KindID, preferenceData.Email)
    }
    return preferences, nil
}
