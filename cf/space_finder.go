package cf

import "github.com/pivotal-cf-experimental/rainmaker"

type tokenGetter interface {
	GetToken(id, secret string) (token string, err error)
}

type spaceGetter interface {
	Get(guid, token string) (rainmaker.Space, error)
}

type SpaceFinder struct {
	clientID     string
	clientSecret string
	clients      tokenGetter
	spaces       spaceGetter
}

func NewSpaceFinder(clientID, clientSecret string, clients tokenGetter, spaces spaceGetter) SpaceFinder {
	return SpaceFinder{
		clientID:     clientID,
		clientSecret: clientSecret,
		clients:      clients,
		spaces:       spaces,
	}
}

func (f SpaceFinder) Exists(guid string) (bool, error) {
	token, err := f.clients.GetToken(f.clientID, f.clientSecret)
	if err != nil {
		return false, err
	}

	_, err = f.spaces.Get(guid, token)
	if err != nil {
		switch err.(type) {
		case rainmaker.NotFoundError:
			return false, nil
		}
		return false, err
	}

	return true, nil
}
