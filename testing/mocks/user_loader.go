package mocks

import "github.com/cloudfoundry-incubator/notifications/uaa"

type UserLoader struct {
	LoadCall struct {
		Receives struct {
			UserGUIDs []string
			Token     string
		}
		Returns struct {
			Users map[string]uaa.User
			Error error
		}
	}
}

func NewUserLoader() *UserLoader {
	return &UserLoader{}
}

func (ul *UserLoader) Load(userGUIDs []string, token string) (map[string]uaa.User, error) {
	ul.LoadCall.Receives.UserGUIDs = userGUIDs
	ul.LoadCall.Receives.Token = token

	return ul.LoadCall.Returns.Users, ul.LoadCall.Returns.Error
}
