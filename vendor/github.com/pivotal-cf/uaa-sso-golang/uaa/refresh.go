package uaa

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

type RefreshInterface interface {
	Refresh(string) (Token, error)
}

func Refresh(u UAA, refreshToken string) (Token, error) {
	token := NewToken()
	params := url.Values{
		"grant_type":    {"refresh_token"},
		"redirect_uri":  {u.RedirectURL},
		"refresh_token": {refreshToken},
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

	switch {
	case code == http.StatusUnauthorized:
		return token, InvalidRefreshToken
	case code > 399:
		return token, NewFailure(code, body)
	}

	json.Unmarshal(body, &token)
	return token, nil
}
