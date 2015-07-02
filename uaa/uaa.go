package uaa

import (
	uaaSSOGolang "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

type UAAClient struct {
	Client *uaaSSOGolang.UAA
}

func NewUAAClient(host, clientID, clientSecret string, verifySSL bool) (client UAAClient) {
	uaaSSOGolangClient := uaaSSOGolang.NewUAA("", host, clientID, clientSecret, "")
	client.Client = &uaaSSOGolangClient
	client.Client.VerifySSL = verifySSL
	return client
}

func (u *UAAClient) SetToken(token string) {
	u.Client.SetToken(token)
}

func (u *UAAClient) GetClientToken() (uaaSSOGolang.Token, error) {
	token, err := u.Client.GetClientToken()
	return token, err
}

func (u *UAAClient) UsersGUIDsByScope(scope string) ([]string, error) {
	guids, err := u.Client.UsersGUIDsByScope(scope)
	return guids, err
}

func (u *UAAClient) AllUsers() ([]uaaSSOGolang.User, error) {
	users, err := u.Client.AllUsers()
	return users, err
}

func (u *UAAClient) UsersEmailsByIDs(ids ...string) ([]uaaSSOGolang.User, error) {
	users, err := u.Client.UsersEmailsByIDs(ids...)
	return users, err
}

func (u *UAAClient) GetTokenKey() (string, error) {
	key, err := u.Client.GetTokenKey()
	return key, err
}
