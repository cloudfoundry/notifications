package fakes

import (
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
)

type PreferencesFinder struct {
	FindCall struct {
		Receives struct {
			Database models.DatabaseInterface
			UserGUID string
		}
		Returns struct {
			PreferencesBuilder services.PreferencesBuilder
			Error              error
		}
	}
}

func NewPreferencesFinder(builder services.PreferencesBuilder) *PreferencesFinder {
	finder := &PreferencesFinder{}
	finder.FindCall.Returns.PreferencesBuilder = builder

	return finder
}

func (pb *PreferencesFinder) Find(database models.DatabaseInterface, userGUID string) (services.PreferencesBuilder, error) {
	pb.FindCall.Receives.Database = database
	pb.FindCall.Receives.UserGUID = userGUID

	return pb.FindCall.Returns.PreferencesBuilder, pb.FindCall.Returns.Error
}
