package fakes

import (
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/web/services"
)

type PreferencesFinder struct {
	ReturnValue services.PreferencesBuilder

	FindCall struct {
		Arguments []interface{}
		Error     error
	}
}

func NewPreferencesFinder(returnValue services.PreferencesBuilder) *PreferencesFinder {
	return &PreferencesFinder{
		ReturnValue: returnValue,
	}
}

func (fake *PreferencesFinder) Find(database models.DatabaseInterface, userGUID string) (services.PreferencesBuilder, error) {
	fake.FindCall.Arguments = []interface{}{database, userGUID}
	return fake.ReturnValue, fake.FindCall.Error
}
