package uaa

import (
	"fmt"
	"net/url"
	"time"

	"github.com/dgrijalva/jwt-go"
	uaaSSOGolang "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

type ZonedUAAClient struct {
	clientID     string
	clientSecret string
	verifySSL    bool
	UAAPublicKey string
}

func NewZonedUAAClient(clientID, clientSecret string, verifySSL bool, uaaPublicKey string) (client ZonedUAAClient) {
	return ZonedUAAClient{
		clientID:     clientID,
		clientSecret: clientSecret,
		verifySSL:    verifySSL,
		UAAPublicKey: uaaPublicKey,
	}
}

func (z ZonedUAAClient) ZonedGetClientToken(host string) (string, error) {
	uaaClient := uaaSSOGolang.NewUAA("", host, z.clientID, z.clientSecret, "")
	uaaClient.VerifySSL = z.verifySSL
	token, err := uaaClient.GetClientToken()
	return token.Access, err
}

func (z ZonedUAAClient) UsersEmailsByIDs(token string, ids ...string) ([]User, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(z.UAAPublicKey), nil
	})
	if err != nil {
		return nil, err
	}

	tokenIssuerURL, err := url.Parse(parsedToken.Claims["iss"].(string))
	if err != nil {
		return nil, err
	}
	uaaHost := tokenIssuerURL.Scheme + "://" + tokenIssuerURL.Host
	uaaClient := uaaSSOGolang.NewUAA("", uaaHost, z.clientID, z.clientSecret, "")
	uaaClient.VerifySSL = z.verifySSL
	uaaClient.SetToken(token)

	users, err := uaaClient.UsersEmailsByIDs(ids...)
	myUsers := make([]User, len(users))
	for index, user := range users {
		myUsers[index].fromSSOGolangUser(user)
	}
	return myUsers, err
}

type UAAClient struct {
	Client *uaaSSOGolang.UAA
}

type User struct {
	ID     string
	Emails []string
}

func NewUAAClient(host, clientID, clientSecret string, verifySSL bool) (client UAAClient) {
	uaaSSOGolangClient := uaaSSOGolang.NewUAA("", host, clientID, clientSecret, "")
	client.Client = &uaaSSOGolangClient
	client.Client.VerifySSL = verifySSL
	return client
}

func (u *UAAClient) SetToken(token string) {
	u.Client.SetToken(token)
}

func (u *UAAClient) GetClientToken() (string, error) {
	token, err := u.Client.GetClientToken()
	return token.Access, err
}

func (u *UAAClient) UsersGUIDsByScope(scope string) ([]string, error) {
	guids, err := u.Client.UsersGUIDsByScope(scope)
	return guids, err
}

func (u *User) fromSSOGolangUser(user uaaSSOGolang.User) {
	u.ID = user.ID
	u.Emails = user.Emails
}

func (u *UAAClient) AllUsers() ([]User, error) {
	users, err := u.Client.AllUsers()

	myUsers := make([]User, len(users))
	for index, user := range users {
		myUsers[index].fromSSOGolangUser(user)
	}
	return myUsers, err
}

func (u *UAAClient) GetTokenKey() (string, error) {
	key, err := u.Client.GetTokenKey()
	return key, err
}

type Failure struct {
	code    int
	message string
}

func NewFailure(code int, message []byte) Failure {
	return Failure{
		code:    code,
		message: string(message),
	}
}

func (failure Failure) Code() int {
	return failure.code
}

func (failure Failure) Message() string {
	return failure.message
}

func (failure Failure) Error() string {
	return fmt.Sprintf("UAA Wrapper Failure: %d %s", failure.code, failure.message)
}

func AccessTokenExpiresBefore(accessToken string, duration time.Duration) (bool, error) {
	token := uaaSSOGolang.Token{
		Access: accessToken,
	}
	expired, err := token.ExpiresBefore(duration)
	return expired, err
}
