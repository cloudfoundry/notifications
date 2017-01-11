package cf

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/metrics"
)

func (cc CloudController) GetUsersBySpaceGuid(guid, token string) ([]CloudControllerUser, error) {
	then := time.Now()

	list, err := cc.client.Spaces.ListUsers(guid, token)
	if err != nil {
		return []CloudControllerUser{}, NewFailure(0, err.Error())
	}

	duration := time.Now().Sub(then)

	metrics.NewMetric("histogram", map[string]interface{}{
		"name":  "notifications.external-requests.cc.users-by-space-guid",
		"value": duration.Seconds(),
	}).Log()

	ccUsers := []CloudControllerUser{}
	for _, user := range list.Users {
		ccUsers = append(ccUsers, CloudControllerUser{
			GUID: user.GUID,
		})
	}

	return ccUsers, nil
}
