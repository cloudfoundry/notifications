package uaa

import (
	"encoding/json"
	"net/url"
)

type GetTokenKeyInterface interface {
	GetTokenKey() (string, error)
}

func GetTokenKey(u UAA) (string, error) {
	tokenURL := u.uaaURL + "/token_key"
	uri, err := url.Parse(tokenURL)
	if err != nil {
		return "", err
	}

	token, err := GetClientToken(u)
	if err != nil {
		return "", err
	}

	host := uri.Scheme + "://" + uri.Host

	client := NewClient(host, u.VerifySSL).WithAuthorizationToken(token.Access)
	code, body, err := client.MakeRequest("GET", uri.RequestURI(), nil)
	if err != nil {
		return "", err
	}

	if code > 399 {
		return "", NewFailure(code, body)
	}

	hash := make(map[string]interface{})
	json.Unmarshal(body, &hash)
	return hash["value"].(string), nil
}
