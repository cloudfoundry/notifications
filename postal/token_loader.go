package postal

import (
    "sync"
    "time"

    "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

var token uaa.Token
var mutex sync.Mutex

type TokenLoader struct {
    uaaClient UAAInterface
}

func NewTokenLoader(uaaClient UAAInterface) TokenLoader {
    return TokenLoader{
        uaaClient: uaaClient,
    }
}

func (loader TokenLoader) Load() (string, error) {
    var err error

    mutex.Lock()
    defer mutex.Unlock()

    if loader.newTokenRequired() {

        token, err = loader.uaaClient.GetClientToken()
        if err != nil {
            err = UAAErrorFor(err)
            return "", err
        }
    }
    loader.uaaClient.SetToken(token.Access)
    return token.Access, nil
}

func ResetLoader() {
    token = uaa.Token{}
}

func (loader TokenLoader) expired() bool {
    toleratedRequestTime := time.Duration(30 * time.Second)
    expired, _ := token.ExpiresBefore(time.Duration(toleratedRequestTime))
    return expired
}

func (loader TokenLoader) newTokenRequired() bool {
    return token.Access == "" || loader.expired()
}
