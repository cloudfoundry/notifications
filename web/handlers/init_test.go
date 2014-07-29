package handlers_test

import (
    "errors"
    "testing"

    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/mail"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/dgrijalva/jwt-go"
    "github.com/nu7hatch/gouuid"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

func TestWebHandlersSuite(t *testing.T) {
    RegisterFastTokenSigningMethod()

    RegisterFailHandler(Fail)
    RunSpecs(t, "Web Handlers Suite")
}

const (
    UAAPrivateKey = "PRIVATE-KEY"
    UAAPublicKey  = "PUBLIC-KEY"
)

type SigningMethodFast struct{}

func (m SigningMethodFast) Alg() string {
    return "FAST"
}

func (m SigningMethodFast) Sign(signingString string, key []byte) (string, error) {
    signature := jwt.EncodeSegment([]byte(signingString + "SUPERFAST"))
    return signature, nil
}

func (m SigningMethodFast) Verify(signingString, signature string, key []byte) (err error) {
    if signature != jwt.EncodeSegment([]byte(signingString+"SUPERFAST")) {
        return errors.New("Signature is invalid")
    }

    return nil
}

func RegisterFastTokenSigningMethod() {
    jwt.RegisterSigningMethod("FAST", func() jwt.SigningMethod {
        return SigningMethodFast{}
    })
}

func BuildToken(header map[string]interface{}, claims map[string]interface{}) string {
    config.UAAPublicKey = UAAPublicKey

    alg := header["alg"].(string)
    signingMethod := jwt.GetSigningMethod(alg)
    token := jwt.New(signingMethod)
    token.Header = header
    token.Claims = claims

    signed, err := token.SignedString([]byte(UAAPrivateKey))
    if err != nil {
        panic(err)
    }

    return signed
}

type FakeMailClient struct {
    messages       []mail.Message
    errorOnSend    bool
    errorOnConnect bool
}

func (fake *FakeMailClient) Connect() error {
    if fake.errorOnConnect {
        return errors.New("BOOM!")
    }
    return nil
}

func (fake *FakeMailClient) Send(msg mail.Message) error {
    err := fake.Connect()
    if err != nil {
        return err
    }

    if fake.errorOnSend {
        return errors.New("BOOM!")
    }

    fake.messages = append(fake.messages, msg)
    return nil
}

type FakeUAAClient struct {
    ClientToken      uaa.Token
    UsersByID        map[string]uaa.User
    ErrorForUserByID error
}

func (fake FakeUAAClient) SetToken(token string) {}

func (fake FakeUAAClient) GetClientToken() (uaa.Token, error) {
    return fake.ClientToken, nil
}

func (fake FakeUAAClient) UsersByIDs(ids ...string) ([]uaa.User, error) {
    users := []uaa.User{}
    for _, id := range ids {
        if user, ok := fake.UsersByID[id]; ok {
            users = append(users, user)
        }
    }

    return users, fake.ErrorForUserByID
}

type FakeCloudController struct {
    UsersBySpaceGuid         map[string][]cf.CloudControllerUser
    CurrentToken             string
    GetUsersBySpaceGuidError error
    Spaces                   map[string]cf.CloudControllerSpace
    Orgs                     map[string]cf.CloudControllerOrganization
}

func NewFakeCloudController() *FakeCloudController {
    return &FakeCloudController{
        UsersBySpaceGuid: make(map[string][]cf.CloudControllerUser),
    }
}

func (fake *FakeCloudController) GetUsersBySpaceGuid(guid, token string) ([]cf.CloudControllerUser, error) {
    fake.CurrentToken = token

    if users, ok := fake.UsersBySpaceGuid[guid]; ok {
        return users, fake.GetUsersBySpaceGuidError
    } else {
        return make([]cf.CloudControllerUser, 0), fake.GetUsersBySpaceGuidError
    }
}

func (fake *FakeCloudController) LoadSpace(guid, token string) (cf.CloudControllerSpace, error) {
    if space, ok := fake.Spaces[guid]; ok {
        return space, nil
    } else {
        return cf.CloudControllerSpace{}, nil
    }
}

func (fake *FakeCloudController) LoadOrganization(guid, token string) (cf.CloudControllerOrganization, error) {
    if org, ok := fake.Orgs[guid]; ok {
        return org, nil
    } else {
        return cf.CloudControllerOrganization{}, nil
    }
}

var FakeGuidGenerator = postal.GUIDGenerationFunc(func() (*uuid.UUID, error) {
    guid := uuid.UUID([16]byte{0xDE, 0xAD, 0xBE, 0xEF, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55})
    return &guid, nil
})
