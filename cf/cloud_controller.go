package cf

import (
    "bytes"
    "crypto/tls"
    "fmt"
    "io"
    "net/http"
    "sync"

    "github.com/cloudfoundry-incubator/notifications/config"
)

var _client *http.Client
var mutex sync.Mutex

func GetClient() *http.Client {
    mutex.Lock()
    defer mutex.Unlock()

    if _client == nil {
        env := config.NewEnvironment()

        _client = &http.Client{
            Transport: &http.Transport{
                TLSClientConfig: &tls.Config{
                    InsecureSkipVerify: !env.VerifySSL,
                },
            },
        }
    }

    return _client
}

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
    httpClient := GetClient()
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
