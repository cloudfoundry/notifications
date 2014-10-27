package fakes

import "github.com/cloudfoundry-incubator/notifications/web/services"

type FakePreferencesFinder struct {
    ReturnValue services.PreferencesBuilder
    FindError   error
    UserGUID    string
}

func NewFakePreferencesFinder(returnValue services.PreferencesBuilder) *FakePreferencesFinder {
    return &FakePreferencesFinder{
        ReturnValue: returnValue,
    }
}

func (fake *FakePreferencesFinder) Find(userGUID string) (services.PreferencesBuilder, error) {
    fake.UserGUID = userGUID
    return fake.ReturnValue, fake.FindError
}
