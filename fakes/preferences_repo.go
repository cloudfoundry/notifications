package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type FakePreferencesRepo struct {
    NonCriticalPreferences []models.Preference
    FindError              error
}

func NewFakePreferencesRepo(nonCriticalPreferences []models.Preference) *FakePreferencesRepo {
    return &FakePreferencesRepo{
        NonCriticalPreferences: nonCriticalPreferences,
    }
}

func (fake FakePreferencesRepo) FindNonCriticalPreferences(conn models.ConnectionInterface, userGUID string) ([]models.Preference, error) {
    return fake.NonCriticalPreferences, fake.FindError
}
