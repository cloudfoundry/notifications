package uaa

import (
	"fmt"
	"net/url"
	"strings"
)

type UsersEmailsByIDsInterface interface {
	UsersEmailsByIDs(...string) ([]User, error)
}

func UsersEmailsByIDs(uaa UAA, ids ...string) ([]User, error) {
	return UsersEmailsByIDsWithMaxLength(uaa, MaxQueryLength, ids...)
}

func UsersEmailsByIDsWithMaxLength(u UAA, length int, ids ...string) ([]User, error) {
	var filters []string
	var uris []string
	users := []User{}

	for _, id := range ids {
		filters = append(filters, fmt.Sprintf(`Id eq "%s"`, id))
	}

	var start = 0
	for i, _ := range filters {
		if len(UsersEmailsQueryURIFromParts(u.uaaURL, filters[start:i+1])) > length {
			uris = append(uris, UsersEmailsQueryURIFromParts(u.uaaURL, filters[start:i]))
			start = i
		}
	}

	uris = append(uris, UsersEmailsQueryURIFromParts(u.uaaURL, filters[start:]))

	for _, uri := range uris {
		usersToAdd, err := UsersFromQuery(u, uri)
		if err != nil {
			return users, err
		}
		users = append(users, usersToAdd...)
	}
	return users, nil
}

func UsersEmailsQueryURIFromParts(host string, filters []string) string {
	return fmt.Sprintf("%s/Users?attributes=emails,id&filter=%s", host, url.QueryEscape(strings.Join(filters, " or ")))
}
