package fakes

import "github.com/cloudfoundry-incubator/notifications/uaa"

type UAAClient struct {
	ClientAccessToken string
	ClientTokenError  error
	AccessToken       string
	AllUsersError     error
	AllUsersData      []uaa.User
}

type ZonedUAAClient struct {
	ErrorForUserByID          error
	UsersByID                 map[string]uaa.User
	ZonedGetClientTokenHost   string
	ZonedToken                string
	AllUsersData              []uaa.User
	AllUsersError             error
	AllUsersToken             string
	UsersGUIDsByScopeResponse map[string][]string
	UsersGUIDsByScopeError    error
}

func NewUAAClient() *UAAClient {
	return &UAAClient{}
}

func (fake *UAAClient) SetToken(token string) {
	fake.AccessToken = token
}

func (fake UAAClient) GetClientToken() (string, error) {
	return fake.ClientAccessToken, fake.ClientTokenError
}

func NewZonedUAAClient() *ZonedUAAClient {
	return &ZonedUAAClient{
		UsersGUIDsByScopeResponse: make(map[string][]string),
	}
}

func (z *ZonedUAAClient) ZonedGetClientToken(host string) (string, error) {
	z.ZonedGetClientTokenHost = host
	return z.ZonedToken, nil
}

func (z ZonedUAAClient) UsersEmailsByIDs(token string, ids ...string) ([]uaa.User, error) {
	users := []uaa.User{}
	for _, id := range ids {
		if user, ok := z.UsersByID[id]; ok {
			users = append(users, user)
		}
	}

	return users, z.ErrorForUserByID
}

func (z *ZonedUAAClient) UsersGUIDsByScope(token, scope string) ([]string, error) {
	return z.UsersGUIDsByScopeResponse[scope], z.UsersGUIDsByScopeError
}

func (z *ZonedUAAClient) AllUsers(token string) ([]uaa.User, error) {
	z.AllUsersToken = token
	return z.AllUsersData, z.AllUsersError
}
