package cf

import (
	"fmt"
	"time"

	"github.com/pivotal-cf-experimental/rainmaker"
	"github.com/rcrowley/go-metrics"
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

	metrics.GetOrRegisterTimer("notifications.external-requests.cc.space", nil).Update(time.Since(then))

	return CloudControllerSpace{
		GUID:             space.GUID,
		Name:             space.Name,
		OrganizationGUID: space.OrganizationGUID,
	}, nil
}
