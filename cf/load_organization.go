package cf

import (
    "encoding/json"
    "errors"
    "fmt"
)

type CloudControllerOrganization struct {
    Guid string
    Name string
}

type CloudControllerOrganizationResponse struct {
    Metadata struct {
        Guid string `json:"guid"`
    }   `json:"metadata"`
    Entity struct {
        Name string `json:"name"`
    }   `json:"entity"`
}

func (cc CloudController) LoadOrganization(guid, token string) (CloudControllerOrganization, error) {
    client := NewCloudControllerClient(cc.Host)
    org := CloudControllerOrganization{}

    code, body, err := client.MakeRequest("GET", cc.OrganizationPath(guid), token, nil)
    if err != nil {
        return org, err
    }

    if code > 399 {
        return org, errors.New(fmt.Sprintf("CloudController Failure (%d): %s", code, body))
    }

    response := CloudControllerOrganizationResponse{}
    err = json.Unmarshal(body, &response)
    if err != nil {
        return org, err
    }
    org.Guid = response.Metadata.Guid
    org.Name = response.Entity.Name

    return org, nil
}

func (cc CloudController) OrganizationPath(guid string) string {
    return fmt.Sprintf("/v2/organizations/%s", guid)
}
