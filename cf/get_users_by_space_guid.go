package cf

import (
    "encoding/json"
    "fmt"
    "time"

    "github.com/cloudfoundry-incubator/notifications/metrics"
)

type CloudControllerUser struct {
    Guid string `json:"guid"`
}

type CloudControllerUsersResponse struct {
    Resources []struct {
        Metadata struct {
            Guid string `json:"guid"`
        } `json:"metadata"`
    } `json:"resources"`
}

func (cc CloudController) GetUsersBySpaceGuid(guid, token string) ([]CloudControllerUser, error) {
    users := make([]CloudControllerUser, 0)
    client := NewCloudControllerClient(cc.Host)

    then := time.Now()

    code, body, err := client.MakeRequest("GET", cc.UsersBySpaceGuidPath(guid), token, nil)
    if err != nil {
        return users, err
    }

    duration := time.Now().Sub(then)

    metrics.NewMetric("histogram", map[string]interface{}{
        "name":  "notifications.external-requests.cc.users-by-space-guid",
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
            Guid: resource.Metadata.Guid,
        }
        users = append(users, user)
    }

    return users, nil
}

func (cc CloudController) UsersBySpaceGuidPath(guid string) string {
    return fmt.Sprintf("/v2/users?q=space_guid:%s", guid)
}
