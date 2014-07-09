package handlers_test

import (
    "errors"
    "testing"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/mail"
    "github.com/dgrijalva/jwt-go"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

const (
    UAAPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAmH9VDot2sa5NNbs0a7gKwSIN3Pj4r2Q+pcYW8LwcncUTeECt
B6tn3Mz9NvSsa7kdCG2DgVyWM36ru5pEVeP24WHhrc/1cbJjDryRYBAPS+4JJ5r6
NS8B1Rxm7K34BDDQweeQLDtf4d/wG6cjOuvU9Pn3+ePQzwRiC8tyNiybfgsRS3sh
DyqjZluuWHkKrGoxrdtq8N1aGu8C+r5K2lkXKhSISnpxrq3edBVGUXRtbGZ+dXaN
qxKC4bemAokY/+Q8Ip9AnwPdm02Y+XkGwVduVq+5Q2gTtkFgRSUPZyUm27pVAeDC
yBRxYJzUhQjK8HWjHlchuejtsbld9SrfSxQWTwIDAQABAoIBAFZjW8Xfj5/cQ/T6
Vhnnqn/6UKwrhoWlXi/+5aP+jJ97syneSaccnLvijFeDh+GGfkH1+BdiYdxOF+8w
1yFpAMRw9K3ILxz3l1IT1K78qg2zjRAYpUFXncwiSNQvQV7uYHRYP74u7IRCnfys
VDLewkb9DFNNkU6VBw3zdIHoBzYBIyctyMMP45K5ykz00L9ck2HosCqiN4glUaR/
5t65sc2mLZ7Nrpf9rJXQb6khJifYk3DW8UJq6c0BV5eQhMNbv+lWLHwP2c8VStlQ
WfSo1w3Hu4a9QGM/MLZ76j+f9518hnmjXHoZruDOy34PBM3A+8DLgSVv1SCFoMyf
XNnSlfECgYEAySzLmhZWMK//5Ms10cYAHD4UfEiGyPNGGy6hxhVmvGlng6YAe/4b
xmySR6JDLEG7t6IsLaEelW29B0KUfcFr0HYSzWqHJpTzppnlLRBWRnoBmpr7lNe3
QV4JfHDc8AaodQ9OcKdmD2hgb42Mfq53pmZTqGrknFXWzSSaFoxpOrMCgYEAwg59
G7NjaejWaJ/YUURuc0e6VzX7SLkGBdzXxoUXT3ZicS8GOxn49L31OnF1eBpHS1F8
+SLC7VRZP7z9XIKbiYYV1gBQYM/F9RYCPNAxqQJhBQ8M1ZU3Y7iydHk9HiaBBAX5
+OQADy3+uDPbxCLzSfMAlAPdF7lzbgE7gUk98/UCgYEAgrL5rCgq4wLVS33Cf4EV
/UNP59buyotS1sIbFCg/UNViDSPCWMwkm2taNfPzlEM4g/t2nEZ7KjXbg2X8Nx98
vjiXyqEVITnQekKto/NjOfJ2LE3YeUEUrAE+RHzG7aJFu5ewLHx1UDlNvevGhV8w
GQmN/HNGB1O1dB39hfy/OQUCgYAwFwkYAUekqmff+6TO1ueMN/1MuXrxVbDRaR4r
+zWAorTYma+wm8ofVKfd+NoEjnaWirYuw1eNGvcXHY2oDFHhLdJheyhwJW1IRFD/
oxR7brR+XXFvyI+2bcIDrTvhKeeVCKoe7Nm66UoTef5/R64E6Gx/QcnbpECfxTxq
2Ky6tQKBgAwSU8FWvnuN3gtVw/TCd3QsDTBGOU0Nz87ss3yZ3bZdjXmOkMDwxpZx
VDVG0sVM5aSVQHc5B1TQNSQDxfLA1mBrv9AkafwDyLu2Wls6brLw1QCYABHW6CHp
3+QF9DF2DPlkHNHomOQb1Fyz5kkq/fSAVkE5SVkyx3UPjOe2TGoX
-----END RSA PRIVATE KEY-----`
    UAAPublicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAmH9VDot2sa5NNbs0a7gK
wSIN3Pj4r2Q+pcYW8LwcncUTeECtB6tn3Mz9NvSsa7kdCG2DgVyWM36ru5pEVeP2
4WHhrc/1cbJjDryRYBAPS+4JJ5r6NS8B1Rxm7K34BDDQweeQLDtf4d/wG6cjOuvU
9Pn3+ePQzwRiC8tyNiybfgsRS3shDyqjZluuWHkKrGoxrdtq8N1aGu8C+r5K2lkX
KhSISnpxrq3edBVGUXRtbGZ+dXaNqxKC4bemAokY/+Q8Ip9AnwPdm02Y+XkGwVdu
Vq+5Q2gTtkFgRSUPZyUm27pVAeDCyBRxYJzUhQjK8HWjHlchuejtsbld9SrfSxQW
TwIDAQAB
-----END PUBLIC KEY-----`
)

func BuildToken(header map[string]interface{}, claims map[string]interface{}) string {
    config.UAAPublicKey = UAAPublicKey

    token := jwt.New(&jwt.SigningMethodRS256{})
    token.Header = header
    token.Claims = claims

    signed, err := token.SignedString([]byte(UAAPrivateKey))
    if err != nil {
        panic(err)
    }

    return signed
}

func TestWebHandlersSuite(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Web Handlers Suite")
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
    UsersByID        map[string]uaa.User
    ErrorForUserByID error
}

func (fake FakeUAAClient) AuthorizeURL() string {
    return ""
}

func (fake FakeUAAClient) LoginURL() string {
    return ""
}

func (fake FakeUAAClient) SetToken(token string) {}

func (fake FakeUAAClient) Exchange(code string) (uaa.Token, error) {
    return uaa.Token{}, nil
}

func (fake FakeUAAClient) Refresh(token string) (uaa.Token, error) {
    return uaa.Token{}, nil
}

func (fake FakeUAAClient) GetClientToken() (uaa.Token, error) {
    return uaa.Token{}, nil
}

func (fake FakeUAAClient) GetTokenKey() (string, error) {
    return "", nil
}

func (fake FakeUAAClient) UserByID(id string) (uaa.User, error) {
    return fake.UsersByID[id], fake.ErrorForUserByID
}
