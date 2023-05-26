package uaa

import (
	"encoding/json"
	"net/url"
	"strings"
)

type ExchangeInterface interface {
	Exchange(string) (Token, error)
}

func Exchange(u UAA, authCode string) (Token, error) {
	token := NewToken()

	params := url.Values{
		"grant_type":   {"authorization_code"},
		"redirect_uri": {u.RedirectURL},
		"scope":        {u.Scope},
		"code":         {authCode},
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
