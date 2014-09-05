package services

import "github.com/cloudfoundry-incubator/notifications/models"

type PreferencesFinder struct {
    repo models.PreferencesRepoInterface
}

type PreferencesFinderInterface interface {
    Find(string) (PreferencesBuilder, error)
}

func NewPreferencesFinder(repo models.PreferencesRepoInterface) *PreferencesFinder {
    return &PreferencesFinder{
        repo: repo,
    }
}

func (preference PreferencesFinder) Find(userGUID string) (PreferencesBuilder, error) {
    preferences, err := preference.repo.FindNonCriticalPreferences(models.Database().Connection(), userGUID)
    if err != nil {
        return PreferencesBuilder{}, err
    }

    builder := NewPreferencesBuilder()
    for _, preference := range preferences {
        builder.Add(preference)
    }
    return builder, nil
}
