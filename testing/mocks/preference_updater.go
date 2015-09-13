package mocks

import (
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
)

type PreferenceUpdater struct {
	UpdateCall struct {
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
	pu.UpdateCall.Receives.Connection = conn
	pu.UpdateCall.Receives.Preferences = preferences
	pu.UpdateCall.Receives.GlobalUnsubscribe = globalUnsubscribe
	pu.UpdateCall.Receives.UserID = userID

	return pu.UpdateCall.Returns.Error
}
