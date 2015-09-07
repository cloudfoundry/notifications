package mocks

import "github.com/cloudfoundry-incubator/notifications/v1/services"

type PreferencesFinder struct {
	FindCall struct {
		Receives struct {
			Database services.DatabaseInterface
			UserGUID string
		}
		Returns struct {
			PreferencesBuilder services.PreferencesBuilder
			Error              error
		}
	}
}

func NewPreferencesFinder() *PreferencesFinder {
	return &PreferencesFinder{}
}

func (pb *PreferencesFinder) Find(database services.DatabaseInterface, userGUID string) (services.PreferencesBuilder, error) {
	pb.FindCall.Receives.Database = database
	pb.FindCall.Receives.UserGUID = userGUID

	return pb.FindCall.Returns.PreferencesBuilder, pb.FindCall.Returns.Error
}
