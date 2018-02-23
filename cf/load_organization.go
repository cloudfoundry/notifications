package cf

import (
	"fmt"
	"time"

	"github.com/pivotal-cf-experimental/rainmaker"
	metrics "github.com/rcrowley/go-metrics"
)

func (cc CloudController) LoadOrganization(guid, token string) (CloudControllerOrganization, error) {
	then := time.Now()

	org, err := cc.client.Organizations.Get(guid, token)
	if err != nil {
		_, ok := err.(rainmaker.NotFoundError)
		if ok {
			return CloudControllerOrganization{}, NotFoundError{fmt.Sprintf("Organization %q could not be found", guid)}
		} else {
			return CloudControllerOrganization{}, NewFailure(0, err.Error())
		}
	}

	metrics.GetOrRegisterTimer("notifications.external-requests.cc.organization", nil).Update(time.Since(then))

	return CloudControllerOrganization{
		GUID: org.GUID,
		Name: org.Name,
	}, nil
}
