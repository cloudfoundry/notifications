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
    preferences, err := preference.repo.FindNonCriticalPreferences(models.Database().Connection(), UserGUID)
    if err != nil {
        return PreferencesBuilder{}, err
    }

    builder := NewPreferencesBuilder()
    for _, preference := range preferences {
        builder.Add(preference)
    }
    return builder, nil
}
