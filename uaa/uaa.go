package uaa

import (
	"fmt"
	"net/url"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pivotal-cf-experimental/warrant"
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

func (z ZonedUAAClient) GetClientToken(host string) (string, error) {
	uaaClient := warrant.New(warrant.Config{
		Host:          host,
		SkipVerifySSL: !z.verifySSL,
	})
	return uaaClient.Clients.GetToken(z.clientID, z.clientSecret)
}

func (z ZonedUAAClient) UsersEmailsByIDs(token string, ids ...string) ([]User, error) {
	uaaHost, err := z.tokenHost(token)
	if err != nil {
		return nil, err
	}
	uaaClient := warrant.New(warrant.Config{
		Host:          uaaHost,
		SkipVerifySSL: !z.verifySSL,
	})
	myUsers := make([]User, 0, len(ids))
	for _, id := range ids {
		users, err := uaaClient.Users.List(warrant.Query{Filter: fmt.Sprintf("Id eq \"%s\"", id)}, token)
		if err != nil {
			return nil, err
		}
		user := User{}
		for _, warrantUser := range users {
			user.fromWarrantUser(warrantUser)
			myUsers = append(myUsers, user)
		}
	}
	return myUsers, nil
}

func (z ZonedUAAClient) tokenHost(token string) (string, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(z.UAAPublicKey), nil
	})
	if err != nil {
		return "", err
	}

	tokenIssuerURL, err := url.Parse(parsedToken.Claims["iss"].(string))
	if err != nil {
		return "", err
	}
	uaaHost := tokenIssuerURL.Scheme + "://" + tokenIssuerURL.Host
	return uaaHost, nil
}

func (z ZonedUAAClient) AllUsers(token string) ([]User, error) {
	uaaHost, err := z.tokenHost(token)
	if err != nil {
		return nil, err
	}
	uaaSSOGolangClient := uaaSSOGolang.NewUAA("", uaaHost, z.clientID, z.clientSecret, "")
	uaaSSOGolangClient.VerifySSL = z.verifySSL
	users, err := uaaSSOGolangClient.AllUsers()

	myUsers := make([]User, len(users))
	for index, user := range users {
		myUsers[index].fromSSOGolangUser(user)
	}
	return myUsers, err
}

func (z ZonedUAAClient) UsersGUIDsByScope(token string, scope string) ([]string, error) {
	uaaHost, err := z.tokenHost(token)
	if err != nil {
		return nil, err
	}
	uaaSSOGolangClient := uaaSSOGolang.NewUAA("", uaaHost, z.clientID, z.clientSecret, "")
	uaaSSOGolangClient.VerifySSL = z.verifySSL
	guids, err := uaaSSOGolangClient.UsersGUIDsByScope(scope)
	return guids, err
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

func (u *User) fromWarrantUser(user warrant.User) {
	u.ID = user.ID
	u.Emails = user.Emails
}

func (u *User) fromSSOGolangUser(user uaaSSOGolang.User) {
	u.ID = user.ID
	u.Emails = user.Emails
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
