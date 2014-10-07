package fakes

import (
    "errors"

    "github.com/cloudfoundry-incubator/notifications/web/services"
)

type FakePreferencesFinder struct {
    ReturnValue services.PreferencesBuilder
    FindErrors  bool
    UserGUID    string
}

func NewFakePreferencesFinder(returnValue services.PreferencesBuilder) *FakePreferencesFinder {
    return &FakePreferencesFinder{
        ReturnValue: returnValue,
    }
}

func (fake *FakePreferencesFinder) Find(userGUID string) (services.PreferencesBuilder, error) {
    fake.UserGUID = userGUID
    if fake.FindErrors {
        return fake.ReturnValue, errors.New("Meltdown")
    }
    return fake.ReturnValue, nil
}
