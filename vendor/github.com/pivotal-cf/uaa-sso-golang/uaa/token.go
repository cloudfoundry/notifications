package uaa

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	TokenDecodeError = errors.New("Failed to decode token")
	JSONParseError   = errors.New("Failed to parse JSON")
)

// Encapsulates the access and refresh tokens from UAA
type Token struct {
	Access  string `json:"access_token"`
	Refresh string `json:"refresh_token"`
}

func NewToken() Token {
	return Token{}
}

func (token Token) Type() string {
	return "bearer"
}

// Determines if all the token's information is present
func (token Token) IsPresent() bool {
	return token.Access != "" && token.Refresh != ""
}

// Determines if the token has expired
func (token Token) IsExpired() (bool, error) {
	return token.ExpiresBefore(time.Duration(0))
}

// Determines if the token expires by the current time plus the time buffer
func (token Token) ExpiresBefore(timeBuffer time.Duration) (bool, error) {
	parts := strings.Split(token.Access, ".")
	decodedToken, err := jwt.DecodeSegment(parts[1])
	if err != nil {
		return false, TokenDecodeError
	}

	parsedJson := make(map[string]interface{})
	err = json.Unmarshal(decodedToken, &parsedJson)
	if err != nil {
		return false, JSONParseError
	}

	tokenExpiration := parsedJson["exp"].(float64)

	bufferedExpiration := time.Unix(int64(tokenExpiration), 0).Add(timeBuffer)

	return bufferedExpiration.Before(time.Now()), nil
}
