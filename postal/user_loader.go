package postal

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"
)

type uaaUsersInterface interface {
	uaa.UsersEmailsByIDsInterface
	uaa.SetTokenInterface
}

type UserLoader struct {
	uaaClient uaaUsersInterface
}

type UserLoaderInterface interface {
	Load([]string, string) (map[string]uaa.User, error)
}

func NewUserLoader(uaaClient uaaUsersInterface) UserLoader {
	return UserLoader{
		uaaClient: uaaClient,
	}
}

func (loader UserLoader) Load(guids []string, token string) (map[string]uaa.User, error) {
	users := make(map[string]uaa.User)

	loader.uaaClient.SetToken(token)

	usersByIDs, err := loader.fetchUsersByIDs(guids)
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

func (loader UserLoader) fetchUsersByIDs(guids []string) ([]uaa.User, error) {
	then := time.Now()

	usersByIDs, err := loader.uaaClient.UsersEmailsByIDs(guids...)

	duration := time.Now().Sub(then)

	metrics.NewMetric("histogram", map[string]interface{}{
		"name":  "notifications.external-requests.uaa.users-email",
		"value": duration.Seconds(),
	}).Log()

	return usersByIDs, err
}
