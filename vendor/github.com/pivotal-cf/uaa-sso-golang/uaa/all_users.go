package uaa

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type AllUsersInterface interface {
	AllUsers() ([]User, error)
}

func AllUsers(u UAA) ([]User, error) {
	var users []User
	var totalResults int
	var err error

	bareUsersURL := u.uaaURL + "/Users"
	users, totalResults, err = PaginatedUsersFromQuery(u, bareUsersURL)
	if err != nil {
		return users, err
	}

	for ThereAreMorePages(users, totalResults) {
		var moreUsers []User

		nextStartIndex := len(users) + 1
		moreUsers, totalResults, err = PaginatedUsersFromQuery(u, UsersQueryURIFromStartIndex(u.uaaURL, nextStartIndex))
		if err != nil {
			return users, err
		}
		users = append(users, moreUsers...)
	}

	return users, nil
}

func PaginatedUsersFromQuery(u UAA, uriString string) ([]User, int, error) {
	users := []User{}
	uri, err := url.Parse(uriString)
	if err != nil {
		return []User{}, 0, err
	}

	host := uri.Scheme + "://" + uri.Host
	client := NewClient(host, u.VerifySSL).WithAuthorizationToken(u.AccessToken)
	code, body, err := client.MakeRequest("GET", uri.RequestURI(), nil)
	if err != nil {
		return users, 0, err
	}

	if code > 399 {
		return users, 0, NewFailure(code, body)
	}

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return users, 0, err
	}

	resources := response["resources"].([]interface{})
	for _, resource := range resources {
		user, err := UserFromResource(resource.(map[string]interface{}))
		if err != nil {
			return users, 0, err
		}
		users = append(users, user)
	}

	totalResults := int(response["totalResults"].(float64))

	return users, totalResults, nil
}

func UsersQueryURIFromStartIndex(host string, startIndex int) string {
	return fmt.Sprintf("%s/Users?startIndex=%d", host, startIndex)
}

func ThereAreMorePages(users []User, totalResults int) bool {
	return (len(users) < totalResults)
}
