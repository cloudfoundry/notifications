package cf

import (
    "encoding/json"
    "time"

    "github.com/cloudfoundry-incubator/notifications/metrics"
)

type CloudControllerSpace struct {
    GUID             string
    Name             string
    OrganizationGUID string
}

type CloudControllerSpaceResponse struct {
    Metadata struct {
        GUID string `json:"guid"`
    }   `json:"metadata"`
    Entity struct {
        Name             string `json:"name"`
        OrganizationGUID string `json:"organization_guid"`
    }   `json:"entity"`
}

func (cc CloudController) LoadSpace(spaceGuid, token string) (CloudControllerSpace, error) {
    space := CloudControllerSpace{}

    then := time.Now()

    code, body, err := cc.client.MakeRequest("GET", cc.SpacePath(spaceGuid), token, nil)
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
    space.GUID = spaceResponse.Metadata.GUID
    space.Name = spaceResponse.Entity.Name
    space.OrganizationGUID = spaceResponse.Entity.OrganizationGUID

    return space, nil
}

func (cc CloudController) SpacePath(guid string) string {
    return "/v2/spaces/" + guid
}
