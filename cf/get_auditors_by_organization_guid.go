package cf

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cloudfoundry-incubator/notifications/metrics"
)

func (cc CloudController) GetAuditorsByOrgGuid(guid, token string) ([]CloudControllerUser, error) {
	users := make([]CloudControllerUser, 0)

	then := time.Now()

	code, body, err := cc.client.MakeRequest("GET", cc.AuditorsByOrgGuidPath(guid), token, nil)
	if err != nil {
		return users, err
	}

	duration := time.Now().Sub(then)

	metrics.NewMetric("histogram", map[string]interface{}{
		"name":  "notifications.external-requests.cc.auditors-by-org-guid",
		"value": duration.Seconds(),
	}).Log()

	if code > 399 {
		return users, NewFailure(code, string(body))
	}

	usersResponse := CloudControllerUsersResponse{}
	err = json.Unmarshal(body, &usersResponse)
	if err != nil {
		return users, err
	}

	for _, resource := range usersResponse.Resources {
		user := CloudControllerUser{
			GUID: resource.Metadata.GUID,
		}
		users = append(users, user)
	}

	return users, nil
}

func (cc CloudController) AuditorsByOrgGuidPath(guid string) string {
	return fmt.Sprintf("/v2/organizations/%s/auditors", guid)
}
