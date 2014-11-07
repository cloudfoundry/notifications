package fakes

import "github.com/cloudfoundry-incubator/notifications/web/services"

type PreferencesFinder struct {
    ReturnValue services.PreferencesBuilder
    FindError   error
    UserGUID    string
}

func NewPreferencesFinder(returnValue services.PreferencesBuilder) *PreferencesFinder {
    return &PreferencesFinder{
        ReturnValue: returnValue,
    }
}

func (fake *PreferencesFinder) Find(userGUID string) (services.PreferencesBuilder, error) {
    fake.UserGUID = userGUID
    return fake.ReturnValue, fake.FindError
}
