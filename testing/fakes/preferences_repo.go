package fakes

import (
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/models"
)

type PreferencesRepo struct {
	FindNonCriticalPreferencesCall struct {
		Receives struct {
			Connection db.ConnectionInterface
			UserGUID   string
		}
		Returns struct {
			Preferences []models.Preference
			Error       error
		}
	}
}

func NewPreferencesRepo() *PreferencesRepo {
	return &PreferencesRepo{}
}

func (pr *PreferencesRepo) FindNonCriticalPreferences(conn db.ConnectionInterface, userGUID string) ([]models.Preference, error) {
	pr.FindNonCriticalPreferencesCall.Receives.Connection = conn
	pr.FindNonCriticalPreferencesCall.Receives.UserGUID = userGUID

	return pr.FindNonCriticalPreferencesCall.Returns.Preferences, pr.FindNonCriticalPreferencesCall.Returns.Error
}
