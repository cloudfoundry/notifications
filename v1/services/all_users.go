package services

import "github.com/cloudfoundry-incubator/notifications/uaa"

type AllUsers struct {
	uaa uaaAllUsersInterface
}

type uaaAllUsersInterface interface {
	AllUsers(string) ([]uaa.User, error)
}

func NewAllUsers(uaa uaaAllUsersInterface) AllUsers {
	return AllUsers{
		uaa: uaa,
	}
}

func (allUsers AllUsers) AllUserGUIDs(token string) ([]string, error) {
	var guids []string

	usersMap, err := allUsers.uaa.AllUsers(token)
	if err != nil {
		return guids, err
	}

	for _, user := range usersMap {
		guids = append(guids, user.ID)
	}

	return guids, nil
}
