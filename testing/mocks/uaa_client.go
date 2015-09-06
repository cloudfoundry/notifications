package mocks

import "github.com/cloudfoundry-incubator/notifications/uaa"

type ZonedUAAClient struct {
	AllUsersCall struct {
		Receives struct {
			Token string
		}
		Returns struct {
			Users []uaa.User
			Error error
		}
	}

	UsersGUIDsByScopeCall struct {
		Receives struct {
			Token string
			Scope string
		}
		Returns struct {
			UserGUIDs []string
			Error     error
		}
	}

	GetClientTokenCall struct {
		Receives struct {
			Host string
		}
		Returns struct {
			Token string
			Error error
		}
	}

	UsersEmailsByIDsCall struct {
		Receives struct {
			Token string
			IDs   []string
		}
		Returns struct {
			Users []uaa.User
			Error error
		}
	}
}

func NewZonedUAAClient() *ZonedUAAClient {
	return &ZonedUAAClient{}
}

func (c *ZonedUAAClient) AllUsers(token string) ([]uaa.User, error) {
	c.AllUsersCall.Receives.Token = token

	return c.AllUsersCall.Returns.Users, c.AllUsersCall.Returns.Error
}

func (c *ZonedUAAClient) UsersGUIDsByScope(token, scope string) ([]string, error) {
	c.UsersGUIDsByScopeCall.Receives.Token = token
	c.UsersGUIDsByScopeCall.Receives.Scope = scope

	return c.UsersGUIDsByScopeCall.Returns.UserGUIDs, c.UsersGUIDsByScopeCall.Returns.Error
}

func (c *ZonedUAAClient) GetClientToken(host string) (string, error) {
	c.GetClientTokenCall.Receives.Host = host

	return c.GetClientTokenCall.Returns.Token, c.GetClientTokenCall.Returns.Error
}

func (c *ZonedUAAClient) UsersEmailsByIDs(token string, ids ...string) ([]uaa.User, error) {
	c.UsersEmailsByIDsCall.Receives.Token = token
	c.UsersEmailsByIDsCall.Receives.IDs = ids

	return c.UsersEmailsByIDsCall.Returns.Users, c.UsersEmailsByIDsCall.Returns.Error
}
