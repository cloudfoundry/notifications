package common

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/uaa"
	metrics "github.com/rcrowley/go-metrics"
)

type uaaEmailGetter interface {
	UsersEmailsByIDs(string, ...string) ([]uaa.User, error)
}

type UserLoader struct {
	uaaClient uaaEmailGetter
}

func NewUserLoader(uaaClient uaaEmailGetter) UserLoader {
	return UserLoader{
		uaaClient: uaaClient,
	}
}

func (loader UserLoader) Load(guids []string, token string) (map[string]uaa.User, error) {
	users := make(map[string]uaa.User)

	usersByIDs, err := loader.fetchUsersByIDs(token, guids)
	if err != nil {
		err = UAAErrorFor(err)
		return users, err
	}

	for _, user := range usersByIDs {
		users[user.ID] = user
	}

	for _, guid := range guids {
		if _, ok := users[guid]; !ok {
			users[guid] = uaa.User{}
		}
	}

	return users, nil
}

func (loader UserLoader) fetchUsersByIDs(token string, guids []string) ([]uaa.User, error) {
	then := time.Now()

	usersByIDs, err := loader.uaaClient.UsersEmailsByIDs(token, guids...)

	metrics.GetOrRegisterTimer("notifications.external-requests.uaa.users-email", nil).Update(time.Since(then))

	return usersByIDs, err
}
