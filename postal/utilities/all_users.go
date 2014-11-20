package utilities

import "github.com/pivotal-cf/uaa-sso-golang/uaa"

type AllUsersInterface interface {
	AllUserEmailsAndGUIDs() (map[string]uaa.User, []string, error)
}

type AllUsers struct {
	uaa UAAInterface
}

func NewAllUsers(uaa UAAInterface) AllUsers {
	return AllUsers{
		uaa: uaa,
	}
}

func (allUsers AllUsers) AllUserEmailsAndGUIDs() (map[string]uaa.User, []string, error) {
	guids := []string{}
	formattedUsers := make(map[string]uaa.User)

	usersMap, err := allUsers.uaa.AllUsers()
	if err != nil {
		return formattedUsers, guids, err
	}

	for _, user := range usersMap {
		guids = append(guids, user.ID)
		formattedUsers[user.ID] = user
	}

	return formattedUsers, guids, nil
}
