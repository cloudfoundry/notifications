package cf

import (
	"fmt"
	"time"

	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/pivotal-cf-experimental/rainmaker"
)

func (cc CloudController) LoadSpace(spaceGuid, token string) (CloudControllerSpace, error) {
	then := time.Now()

	space, err := cc.client.Spaces.Get(spaceGuid, token)
	if err != nil {
		_, ok := err.(rainmaker.NotFoundError)
		if ok {
			return CloudControllerSpace{}, NotFoundError{fmt.Sprintf("Space %q could not be found", spaceGuid)}
		} else {
			return CloudControllerSpace{}, NewFailure(0, err.Error())
		}
	}

	duration := time.Now().Sub(then)

	metrics.NewMetric("histogram", map[string]interface{}{
		"name":  "notifications.external-requests.cc.space",
		"value": duration.Seconds(),
	}).Log()

	return CloudControllerSpace{
		GUID:             space.GUID,
		Name:             space.Name,
		OrganizationGUID: space.OrganizationGUID,
	}, nil
}
