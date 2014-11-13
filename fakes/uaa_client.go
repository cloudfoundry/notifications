package fakes

import "github.com/pivotal-cf/uaa-sso-golang/uaa"

type UAAClient struct {
	ClientToken               uaa.Token
	ClientTokenError          error
	UsersByID                 map[string]uaa.User
	ErrorForUserByID          error
	AccessToken               string
	UsersGUIDsByScopeResponse map[string][]string
	UsersGUIDsByScopeError    error
}

func NewUAAClient() *UAAClient {
	return &UAAClient{
		UsersGUIDsByScopeResponse: make(map[string][]string),
	}
}

func (fake *UAAClient) SetToken(token string) {
	fake.AccessToken = token
}

func (fake UAAClient) GetClientToken() (uaa.Token, error) {
	return fake.ClientToken, fake.ClientTokenError
}

func (fake UAAClient) UsersEmailsByIDs(ids ...string) ([]uaa.User, error) {
	users := []uaa.User{}
	for _, id := range ids {
		if user, ok := fake.UsersByID[id]; ok {
			users = append(users, user)
		}
	}

	return users, fake.ErrorForUserByID
}

func (fake *UAAClient) UsersGUIDsByScope(scope string) ([]string, error) {
	return fake.UsersGUIDsByScopeResponse[scope], fake.UsersGUIDsByScopeError
}
