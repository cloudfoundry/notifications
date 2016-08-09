package uaa

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type UsersGUIDsByScopeInterface interface {
	UsersGUIDsByScope(string) ([]string, error)
}

type userGUIDsByScopeResponse struct {
	Resources []struct {
		Members []struct {
			GUID string `json:"value"`
		} `json:"members"`
	} `json:"resources"`
}

func UsersGUIDsByScope(u UAA, scope string) ([]string, error) {
	var guids []string
	filterValue := url.QueryEscape("displayName eq \"" + scope + "\"")

	uri, err := url.Parse(fmt.Sprintf("%s/Groups?attributes=members&filter=%s", u.uaaURL, filterValue))
	if err != nil {
		return guids, err
	}

	host := uri.Scheme + "://" + uri.Host
	client := NewClient(host, u.VerifySSL).WithAuthorizationToken(u.AccessToken)
	code, body, err := client.MakeRequest("GET", uri.RequestURI(), nil)
	if err != nil {
		return guids, err
	}

	if code > 399 {
		return guids, NewFailure(code, body)
	}

	guids, err = guidsFromBody(body)
	if err != nil {
		return guids, err
	}

	return guids, nil
}

func guidsFromBody(body []byte) ([]string, error) {
	var response userGUIDsByScopeResponse
	guids := []string{}
	err := json.Unmarshal(body, &response)
	if err != nil {
		return guids, err
	}

	for _, member := range response.Resources[0].Members {
		guids = append(guids, member.GUID)
	}

	return guids, nil
}
