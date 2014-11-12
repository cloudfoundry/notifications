package fakes

import "github.com/pivotal-cf/uaa-sso-golang/uaa"

type UserLoader struct {
	Users       map[string]uaa.User
	LoadError   error
	LoadedGUIDs []string
}

func NewUserLoader() *UserLoader {
	return &UserLoader{
		Users:       make(map[string]uaa.User),
		LoadedGUIDs: make([]string, 0),
	}
}

func (fake *UserLoader) Load(guids []string, token string) (map[string]uaa.User, error) {
	fake.LoadedGUIDs = guids
	return fake.Users, fake.LoadError
}
