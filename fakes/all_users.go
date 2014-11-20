package fakes

import "github.com/pivotal-cf/uaa-sso-golang/uaa"

type AllUsers struct {
	Users     map[string]uaa.User
	LoadError error
	GUIDS     []string
}

func NewAllUsers() *AllUsers {
	return &AllUsers{}
}

func (fake *AllUsers) AllUserEmailsAndGUIDs() (map[string]uaa.User, []string, error) {
	return fake.Users, fake.GUIDS, fake.LoadError
}
