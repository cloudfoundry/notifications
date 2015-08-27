package cf

import "github.com/pivotal-cf-experimental/rainmaker"

type OrgFinder struct {
	orgs         orgGetter
	clients      tokenGetter
	clientID     string
	clientSecret string
}

type orgGetter interface {
	Get(guid, token string) (rainmaker.Organization, error)
}

func NewOrgFinder(clientID, clientSecret string, clients tokenGetter, orgs orgGetter) OrgFinder {
	return OrgFinder{
		clients:      clients,
		orgs:         orgs,
		clientSecret: clientSecret,
		clientID:     clientID,
	}
}

func (f OrgFinder) Exists(guid string) (bool, error) {
	token, err := f.clients.GetToken(f.clientID, f.clientSecret)
	if err != nil {
		return false, err
	}

	_, err = f.orgs.Get(guid, token)
	if err != nil {
		switch err.(type) {
		case rainmaker.NotFoundError:
			return false, nil
		}
		return false, err
	}

	return true, nil
}
