/*
Package to interact with Cloudfoundry UAA server.
Constructors are generally provided for objects a client needs to use


 This link is helpful for understanding UAA OAUTH handshakes: http://blog.cloudfoundry.org/2012/07/23/uaa-intro/
*/
package uaa

import (
	"errors"
	"fmt"
	"net/url"
)

var InvalidRefreshToken = errors.New("UAA Invalid Refresh Token")

// used to encapuslate info about errors
type Failure struct {
	code    int
	message string
}

// Failure constructor
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
	return fmt.Sprintf("UAA Failure: %d %s", failure.code, failure.message)
}

// Defines methods needed for clients to use UAA
type UAAInterface interface {
	AuthorizeURLInterface
	LoginURLInterface
	SetTokenInterface
	ExchangeInterface
	GetClientTokenInterface
	GetTokenKeyInterface
	RefreshInterface
	UserByIDInterface
	UsersByIDsInterface
	UsersEmailsByIDsInterface
	UsersGUIDsByScopeInterface
	AllUsersInterface
}

type AuthorizeURLInterface interface {
	AuthorizeURL() string
}

type LoginURLInterface interface {
	LoginURL() string
}

type SetTokenInterface interface {
	SetToken(string)
}

// Contains necessary info to communicate with Cloudfoundry UAA server, use
// the NewUAA constructor to create one.
type UAA struct {
	loginURL       string
	uaaURL         string
	ClientID       string
	ClientSecret   string
	RedirectURL    string
	Scope          string
	State          string
	AccessType     string
	ApprovalPrompt string
	AccessToken    string
	VerifySSL      bool

	ExchangeCommand          func(UAA, string) (Token, error)
	RefreshCommand           func(UAA, string) (Token, error)
	GetClientTokenCommand    func(UAA) (Token, error)
	UserByIDCommand          func(UAA, string) (User, error)
	GetTokenKeyCommand       func(UAA) (string, error)
	UsersByIDsCommand        func(UAA, ...string) ([]User, error)
	UsersEmailsByIDsCommand  func(UAA, ...string) ([]User, error)
	UsersGUIDsByScopeCommand func(UAA, string) ([]string, error)
	AllUsersCommand          func(UAA) ([]User, error)
}

func NewUAA(loginURL, uaaURL, clientID, clientSecret, token string) UAA {
	return UAA{
		loginURL:                 loginURL,
		uaaURL:                   uaaURL,
		ClientID:                 clientID,
		ClientSecret:             clientSecret,
		AccessToken:              token,
		VerifySSL:                true,
		ExchangeCommand:          Exchange,
		GetClientTokenCommand:    GetClientToken,
		GetTokenKeyCommand:       GetTokenKey,
		RefreshCommand:           Refresh,
		UserByIDCommand:          UserByID,
		UsersByIDsCommand:        UsersByIDs,
		UsersEmailsByIDsCommand:  UsersEmailsByIDs,
		UsersGUIDsByScopeCommand: UsersGUIDsByScope,
		AllUsersCommand:          AllUsers,
	}
}

func (u UAA) AuthorizeURL() string {
	return fmt.Sprintf("%s/oauth/authorize", u.loginURL)
}

// Returns url used to login to UAA
func (u UAA) LoginURL() string {
	v := url.Values{}
	v.Set("access_type", u.AccessType)
	v.Set("approval_prompt", u.ApprovalPrompt)
	v.Set("client_id", u.ClientID)
	v.Set("redirect_uri", u.RedirectURL)
	v.Set("response_type", "code")
	v.Set("scope", u.Scope)
	v.Set("state", u.State)

	return u.AuthorizeURL() + "?" + v.Encode()
}

func (u *UAA) SetToken(token string) {
	u.AccessToken = token
}

func (u UAA) tokenURL() string {
	return fmt.Sprintf("%s/oauth/token", u.uaaURL)
}

// Gets auth token based on the code UAA provides during redirect process
func (u UAA) Exchange(authCode string) (Token, error) {
	return u.ExchangeCommand(u, authCode)
}

// Refreshes token from UAA server
func (u UAA) Refresh(refreshToken string) (Token, error) {
	return u.RefreshCommand(u, refreshToken)
}

// Retrieves ClientToken from UAA server
func (u UAA) GetClientToken() (Token, error) {
	return u.GetClientTokenCommand(u)
}

// Retrieves User from UAA server using the user id
func (u UAA) UserByID(id string) (User, error) {
	return u.UserByIDCommand(u, id)
}

func (u UAA) GetTokenKey() (string, error) {
	return u.GetTokenKeyCommand(u)
}

func (u UAA) UsersByIDs(ids ...string) ([]User, error) {
	return u.UsersByIDsCommand(u, ids...)
}

func (u UAA) UsersEmailsByIDs(ids ...string) ([]User, error) {
	return u.UsersEmailsByIDsCommand(u, ids...)
}

func (u UAA) UsersGUIDsByScope(scope string) ([]string, error) {
	return u.UsersGUIDsByScopeCommand(u, scope)
}

func (u UAA) AllUsers() ([]User, error) {
	return u.AllUsersCommand(u)
}
