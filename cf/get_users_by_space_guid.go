package cf

import (
	"time"

	"github.com/rcrowley/go-metrics"
)

func (cc CloudController) GetUsersBySpaceGuid(guid, token string) ([]CloudControllerUser, error) {
	then := time.Now()

	list, err := cc.client.Spaces.ListUsers(guid, token)
	if err != nil {
		return []CloudControllerUser{}, NewFailure(0, err.Error())
	}

	metrics.GetOrRegisterTimer("notifications.external-requests.cc.users-by-space-guid", nil).Update(time.Since(then))

	ccUsers := []CloudControllerUser{}
	for _, user := range list.Users {
		ccUsers = append(ccUsers, CloudControllerUser{
			GUID: user.GUID,
		})
	}

	return ccUsers, nil
}
