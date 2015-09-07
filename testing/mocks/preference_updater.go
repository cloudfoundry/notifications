package mocks

import (
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
)

type PreferenceUpdater struct {
	ExecuteCall struct {
		Receives struct {
			Connection        services.ConnectionInterface
			Preferences       []models.Preference
			GlobalUnsubscribe bool
			UserID            string
		}
		Returns struct {
			Error error
		}
	}
}

func NewPreferenceUpdater() *PreferenceUpdater {
	return &PreferenceUpdater{}
}

func (pu *PreferenceUpdater) Update(conn services.ConnectionInterface, preferences []models.Preference, globalUnsubscribe bool, userID string) error {
	pu.ExecuteCall.Receives.Connection = conn
	pu.ExecuteCall.Receives.Preferences = preferences
	pu.ExecuteCall.Receives.GlobalUnsubscribe = globalUnsubscribe
	pu.ExecuteCall.Receives.UserID = userID

	return pu.ExecuteCall.Returns.Error
}
