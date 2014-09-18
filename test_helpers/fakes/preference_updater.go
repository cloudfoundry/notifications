package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type FakePreferenceUpdater struct {
    ExecuteArguments []interface{}
    ExecuteError     error
}

func NewFakePreferenceUpdater() *FakePreferenceUpdater {
    return &FakePreferenceUpdater{}
}

func (fake *FakePreferenceUpdater) Execute(conn models.ConnectionInterface, preferences []models.Preference, userID string) error {
    fake.ExecuteArguments = append(fake.ExecuteArguments, preferences, userID)
    return fake.ExecuteError
}
