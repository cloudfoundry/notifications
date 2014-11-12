package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type PreferenceUpdater struct {
	ExecuteArguments []interface{}
	ExecuteError     error
}

func NewPreferenceUpdater() *PreferenceUpdater {
	return &PreferenceUpdater{}
}

func (fake *PreferenceUpdater) Execute(conn models.ConnectionInterface, preferences []models.Preference, globalUnsubscribe bool, userID string) error {
	fake.ExecuteArguments = append(fake.ExecuteArguments, preferences, globalUnsubscribe, userID)
	return fake.ExecuteError
}
