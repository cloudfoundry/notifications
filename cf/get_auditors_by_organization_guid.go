package cf

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/metrics"
)

func (cc CloudController) GetAuditorsByOrgGuid(guid, token string) ([]CloudControllerUser, error) {
	ccUsers := make([]CloudControllerUser, 0)

	then := time.Now()

	list, err := cc.client.Organizations.ListAuditors(guid, token)
	if err != nil {
		return ccUsers, NewFailure(0, err.Error())
	}

	duration := time.Now().Sub(then)

	metrics.NewMetric("histogram", map[string]interface{}{
		"name":  "notifications.external-requests.cc.auditors-by-org-guid",
		"value": duration.Seconds(),
	}).Log()

	for _, user := range list.Users {
		ccUsers = append(ccUsers, CloudControllerUser{
			GUID: user.GUID,
		})
	}
	return ccUsers, nil
}
