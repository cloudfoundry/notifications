package utilities

type AllUsersInterface interface {
	AllUserGUIDs() ([]string, error)
}

type AllUsers struct {
	uaa UAAInterface
}

func NewAllUsers(uaa UAAInterface) AllUsers {
	return AllUsers{
		uaa: uaa,
	}
}

func (allUsers AllUsers) AllUserGUIDs() ([]string, error) {
	var guids []string

	usersMap, err := allUsers.uaa.AllUsers()
	if err != nil {
		return guids, err
	}

	for _, user := range usersMap {
		guids = append(guids, user.ID)
	}

	return guids, nil
}
