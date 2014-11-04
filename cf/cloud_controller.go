package cf

import (
    "bytes"
    "crypto/tls"
    "fmt"
    "io"
    "net/http"
    "sync"
)

var _client *http.Client
var mutex sync.Mutex

func GetClient(ccClient CloudControllerClient) *http.Client {
    mutex.Lock()
    defer mutex.Unlock()

    if _client == nil {
        _client = &http.Client{
            Transport: &http.Transport{
                TLSClientConfig: &tls.Config{
                    InsecureSkipVerify: ccClient.skipVerifySSL,
                },
            },
        }
    }

    return _client
}

type CloudController struct {
    client CloudControllerClient
}

type CloudControllerInterface interface {
    GetUsersBySpaceGuid(string, string) ([]CloudControllerUser, error)
    GetUsersByOrgGuid(string, string) ([]CloudControllerUser, error)
    LoadSpace(string, string) (CloudControllerSpace, error)
    LoadOrganization(string, string) (CloudControllerOrganization, error)
}

func NewCloudController(host string, skipVerifySSL bool) CloudController {
    return CloudController{
        client: NewCloudControllerClient(host, skipVerifySSL),
    }
}

type CloudControllerClient struct {
    host          string
    skipVerifySSL bool
}

func NewCloudControllerClient(host string, skipVerifySSL bool) CloudControllerClient {
    return CloudControllerClient{
        host:          host,
        skipVerifySSL: skipVerifySSL,
    }
}

func (client CloudControllerClient) MakeRequest(method, path, token string, body io.Reader) (int, []byte, error) {
    request, err := http.NewRequest(method, client.host+path, body)
    if err != nil {
        return 0, []byte{}, err
    }
    request.Header.Set("Authorization", "Bearer "+token)

    httpClient := GetClient(client)
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

type Failure struct {
    Code    int
    Message string
}

func NewFailure(code int, message string) Failure {
    return Failure{
        Code:    code,
        Message: message,
    }
}

func (failure Failure) Error() string {
    return fmt.Sprintf("CloudController Failure (%d): %s", failure.Code, failure.Message)
}
