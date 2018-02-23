package cf

import (
	"time"

	"github.com/rcrowley/go-metrics"
)

func (cc CloudController) GetBillingManagersByOrgGuid(guid, token string) ([]CloudControllerUser, error) {
	var ccUsers []CloudControllerUser
	then := time.Now()

	list, err := cc.client.Organizations.ListBillingManagers(guid, token)
	if err != nil {
		return ccUsers, NewFailure(0, err.Error())
	}

	metrics.GetOrRegisterTimer("notifications.external-requests.cc.billing-managers-by-org-guid", nil).Update(time.Since(then))

	for _, user := range list.Users {
		ccUsers = append(ccUsers, CloudControllerUser{
			GUID: user.GUID,
		})
	}

	return ccUsers, nil
}
