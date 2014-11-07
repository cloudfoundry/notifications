package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type PreferencesRepo struct {
    NonCriticalPreferences []models.Preference
    FindError              error
}

func NewPreferencesRepo(nonCriticalPreferences []models.Preference) *PreferencesRepo {
    return &PreferencesRepo{
        NonCriticalPreferences: nonCriticalPreferences,
    }
}

func (fake PreferencesRepo) FindNonCriticalPreferences(conn models.ConnectionInterface, userGUID string) ([]models.Preference, error) {
    return fake.NonCriticalPreferences, fake.FindError
}
