package cf

import (
    "bytes"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net/http"
)

type CloudController struct {
    Host string
}

type CloudControllerInterface interface {
    GetUsersBySpaceGuid(string, string) ([]CloudControllerUser, error)
}

func NewCloudController(host string) CloudController {
    return CloudController{
        Host: host,
    }
}

func (cc CloudController) GetUsersBySpaceGuid(guid, token string) ([]CloudControllerUser, error) {
    users := make([]CloudControllerUser, 0)
    client := NewCloudControllerClient(cc.Host)

    code, body, err := client.MakeRequest("GET", cc.UsersBySpaceGuidPath(guid), token, nil)
    if err != nil {
        return users, err
    }

    if code > 399 {
        return users, errors.New(fmt.Sprintf("CloudController Failure (%d): %s", code, body))
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

type CloudControllerClient struct {
    Host string
}

func NewCloudControllerClient(host string) CloudControllerClient {
    return CloudControllerClient{
        Host: host,
    }
}

func (client CloudControllerClient) MakeRequest(method, path, token string, body io.Reader) (int, []byte, error) {
    httpClient := &http.Client{}
    request, err := http.NewRequest(method, client.Host+path, body)
    if err != nil {
        return 0, []byte{}, err
    }
    request.Header.Set("Authorization", "Bearer "+token)

    response, err := httpClient.Do(request)
    if err != nil {
        return 0, []byte{}, err
    }
    code := response.StatusCode

    buffer := bytes.NewBuffer([]byte{})
    _, err = buffer.ReadFrom(response.Body)
    if err != nil {
        return code, []byte{}, err
    }

    return code, buffer.Bytes(), nil
}
