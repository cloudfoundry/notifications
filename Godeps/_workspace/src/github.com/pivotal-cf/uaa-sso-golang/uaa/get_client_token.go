package uaa

import (
	"encoding/json"
	"net/url"
	"strings"
)

type GetClientTokenInterface interface {
	GetClientToken() (Token, error)
}

// Retrieves ClientToken from UAA server
func GetClientToken(u UAA) (Token, error) {
	token := NewToken()
	params := url.Values{
		"grant_type":   {"client_credentials"},
		"redirect_uri": {u.RedirectURL},
	}

	uri, err := url.Parse(u.tokenURL())
	if err != nil {
		return token, err
	}

	host := uri.Scheme + "://" + uri.Host
	client := NewClient(host, u.VerifySSL).WithBasicAuthCredentials(u.ClientID, u.ClientSecret)
	code, body, err := client.MakeRequest("POST", uri.RequestURI(), strings.NewReader(params.Encode()))
	if err != nil {
		return token, err
	}

	if code > 399 {
		return token, NewFailure(code, body)
	}

	json.Unmarshal(body, &token)
	return token, nil
}
