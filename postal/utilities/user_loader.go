package utilities

import (
	"log"
	"time"

	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"
)

type UserLoader struct {
	uaaClient UAAInterface
	logger    *log.Logger
}

type UserLoaderInterface interface {
	Load([]string, string) (map[string]uaa.User, error)
}

func NewUserLoader(uaaClient UAAInterface, logger *log.Logger) UserLoader {
	return UserLoader{
		uaaClient: uaaClient,
		logger:    logger,
	}
}

func (loader UserLoader) Load(guids []string, token string) (map[string]uaa.User, error) {
	users := make(map[string]uaa.User)

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
