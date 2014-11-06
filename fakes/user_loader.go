package fakes

import (
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

type UserLoader struct {
    Users     map[string]uaa.User
    LoadError error
}

func NewUserLoader() *UserLoader {
    return &UserLoader{
        Users: make(map[string]uaa.User),
    }
}

func (fake *UserLoader) Load(postal.TypedGUID, string) (map[string]uaa.User, error) {
    return fake.Users, fake.LoadError
}
