package cf

import (
    "encoding/json"
    "fmt"
    "time"

    "github.com/cloudfoundry-incubator/notifications/metrics"
)

type CloudControllerOrganization struct {
    GUID string
    Name string
}

type CloudControllerOrganizationResponse struct {
    Metadata struct {
        GUID string `json:"guid"`
    }   `json:"metadata"`
    Entity struct {
        Name string `json:"name"`
    }   `json:"entity"`
}

func (cc CloudController) LoadOrganization(guid, token string) (CloudControllerOrganization, error) {
    org := CloudControllerOrganization{}

    then := time.Now()

    code, body, err := cc.client.MakeRequest("GET", cc.OrganizationPath(guid), token, nil)
    if err != nil {
        return org, err
    }

    duration := time.Now().Sub(then)

    metrics.NewMetric("histogram", map[string]interface{}{
        "name":  "notifications.external-requests.cc.organization",
        "value": duration.Seconds(),
    }).Log()

    if code > 399 {
        return org, NewFailure(code, string(body))
    }

    response := CloudControllerOrganizationResponse{}
    err = json.Unmarshal(body, &response)
    if err != nil {
        return org, err
    }
    org.GUID = response.Metadata.GUID
    org.Name = response.Entity.Name

    return org, nil
}

func (cc CloudController) OrganizationPath(guid string) string {
    return fmt.Sprintf("/v2/organizations/%s", guid)
}
