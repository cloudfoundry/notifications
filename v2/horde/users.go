package horde

import "github.com/pivotal-golang/lager"

type Users struct{}

func NewUsers() Users {
	return Users{}
}

func (u Users) GenerateAudiences(guids []string, logger lager.Logger) ([]Audience, error) {
	var users []User
	for _, guid := range guids {
		users = append(users, User{GUID: guid})
	}

	return []Audience{{
		Users:       users,
		Endorsement: "This message was sent directly to you.",
	}}, nil
}
