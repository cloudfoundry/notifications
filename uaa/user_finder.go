package uaa

import (
	"github.com/pivotal-cf-experimental/warrant"
)

type UserFinder struct {
	ID      string
	Secret  string
	Users   userGetter
	Clients tokenGetter
}

type userGetter interface {
	Get(guid, token string) (warrant.User, error)
}

type tokenGetter interface {
	GetToken(id, secret string) (token string, err error)
}

func NewUserFinder(id, secret string, users userGetter, clients tokenGetter) UserFinder {
	return UserFinder{
		ID:      id,
		Secret:  secret,
		Users:   users,
		Clients: clients,
	}
}

func (u UserFinder) Exists(guid string) (bool, error) {
	token, err := u.Clients.GetToken(u.ID, u.Secret)
	if err != nil {
		return false, err
	}

	_, err = u.Users.Get(guid, token)
	if err != nil {
		switch err.(type) {
		case warrant.NotFoundError:
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
