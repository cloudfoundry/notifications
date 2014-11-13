package uaa

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

const MaxQueryLength = 8000

type UsersByIDsInterface interface {
	UsersByIDs(...string) ([]User, error)
}

func UsersByIDs(u UAA, ids ...string) ([]User, error) {
	return UsersByIDsWithMaxLength(u, MaxQueryLength, ids...)
}

func UsersByIDsWithMaxLength(u UAA, length int, ids ...string) ([]User, error) {
	var filters []string
	var uris []string
	users := []User{}

	for _, id := range ids {
		filters = append(filters, fmt.Sprintf(`Id eq "%s"`, id))
	}

	var start = 0
	for i, _ := range filters {
		if len(UsersQueryURIFromParts(u.uaaURL, filters[start:i+1])) > length {
			uris = append(uris, UsersQueryURIFromParts(u.uaaURL, filters[start:i]))
			start = i
		}
	}

	uris = append(uris, UsersQueryURIFromParts(u.uaaURL, filters[start:]))

	for _, uri := range uris {
		usersToAdd, err := UsersFromQuery(u, uri)
		if err != nil {
			return users, err
		}
		users = append(users, usersToAdd...)
	}
	return users, nil
}

func UsersQueryURIFromParts(host string, filters []string) string {
	return fmt.Sprintf("%s/Users?filter=%s", host, url.QueryEscape(strings.Join(filters, " or ")))
}

func UsersFromQuery(u UAA, uriString string) ([]User, error) {
	users := []User{}
	uri, err := url.Parse(uriString)
	if err != nil {
		return []User{}, err
	}

	host := uri.Scheme + "://" + uri.Host
	client := NewClient(host, u.VerifySSL).WithAuthorizationToken(u.AccessToken)
	code, body, err := client.MakeRequest("GET", uri.RequestURI(), nil)
	if err != nil {
		return users, err
	}

	if code > 399 {
		return users, NewFailure(code, body)
	}

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return users, err
	}

	resources := response["resources"].([]interface{})
	for _, resource := range resources {
		user, err := UserFromResource(resource.(map[string]interface{}))
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}

	return users, nil
}
