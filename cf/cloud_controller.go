package cf

import (
    "bytes"
    "crypto/tls"
    "io"
    "net/http"

    "github.com/cloudfoundry-incubator/notifications/config"
)

type CloudController struct {
    Host string
}

type CloudControllerInterface interface {
    GetUsersBySpaceGuid(string, string) ([]CloudControllerUser, error)
    LoadSpace(string, string) (CloudControllerSpace, error)
    LoadOrganization(string, string) (CloudControllerOrganization, error)
}

func NewCloudController(host string) CloudController {
    return CloudController{
        Host: host,
    }
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
    env := config.NewEnvironment()

    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: !env.VerifySSL},
    }
    httpClient := &http.Client{Transport: tr}
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
