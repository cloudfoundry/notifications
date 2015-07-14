package fakes

import "github.com/cloudfoundry-incubator/notifications/uaa"

type UAAClient struct {
	ClientAccessToken         string
	ClientTokenError          error
	AccessToken               string
	UsersGUIDsByScopeResponse map[string][]string
	UsersGUIDsByScopeError    error
	AllUsersError             error
	AllUsersData              []uaa.User
}

type ZonedUAAClient struct {
	ErrorForUserByID        error
	UsersByID               map[string]uaa.User
	ZonedGetClientTokenHost string
	ZonedToken              string
	AllUsersData            []uaa.User
	AllUsersError           error
	AllUsersToken           string
}

func NewUAAClient() *UAAClient {
	return &UAAClient{
		UsersGUIDsByScopeResponse: make(map[string][]string),
	}
}

func (fake *UAAClient) SetToken(token string) {
	fake.AccessToken = token
}

func (fake UAAClient) GetClientToken() (string, error) {
	return fake.ClientAccessToken, fake.ClientTokenError
}

func (fake *UAAClient) UsersGUIDsByScope(scope string) ([]string, error) {
	return fake.UsersGUIDsByScopeResponse[scope], fake.UsersGUIDsByScopeError
}

func NewZonedUAAClient() *ZonedUAAClient {
	return &ZonedUAAClient{}
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

func (z *ZonedUAAClient) AllUsers(token string) ([]uaa.User, error) {
	z.AllUsersToken = token
	return z.AllUsersData, z.AllUsersError
}
