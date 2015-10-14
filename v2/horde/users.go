package horde

type Users struct{}

func NewUsers() Users {
	return Users{}
}

func (u Users) GenerateAudiences(guids []string) ([]Audience, error) {
	var users []User
	for _, guid := range guids {
		users = append(users, User{GUID: guid})
	}

	return []Audience{{
		Users:       users,
		Endorsement: "This message was sent directly to you.",
	}}, nil
}
