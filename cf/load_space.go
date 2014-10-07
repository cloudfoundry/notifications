package cf

import (
    "encoding/json"
    "time"

    "github.com/cloudfoundry-incubator/notifications/metrics"
)

type CloudControllerSpace struct {
    Guid             string
    Name             string
    OrganizationGuid string
}

type CloudControllerSpaceResponse struct {
    Metadata struct {
        Guid string `json:"guid"`
    } `json:"metadata"`
    Entity struct {
        Name             string `json:"name"`
        OrganizationGuid string `json:"organization_guid"`
    } `json:"entity"`
}

func (cc CloudController) LoadSpace(spaceGuid, token string) (CloudControllerSpace, error) {
    client := NewCloudControllerClient(cc.Host)
    space := CloudControllerSpace{}

    then := time.Now()

    code, body, err := client.MakeRequest("GET", cc.SpacePath(spaceGuid), token, nil)
    if err != nil {
        return space, err
    }

    duration := time.Now().Sub(then)

    metrics.NewMetric("histogram", map[string]interface{}{
        "name":  "notifications.external-requests.cc.space",
        "value": duration.Seconds(),
    }).Log()

    if code > 399 {
        return space, NewFailure(code, string(body))
    }

    spaceResponse := CloudControllerSpaceResponse{}
    err = json.Unmarshal(body, &spaceResponse)
    if err != nil {
        return space, err
    }
    space.Guid = spaceResponse.Metadata.Guid
    space.Name = spaceResponse.Entity.Name
    space.OrganizationGuid = spaceResponse.Entity.OrganizationGuid

    return space, nil
}

func (cc CloudController) SpacePath(guid string) string {
    return "/v2/spaces/" + guid
}
