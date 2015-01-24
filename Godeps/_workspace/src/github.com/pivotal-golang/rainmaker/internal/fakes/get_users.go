package fakes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func (fake *CloudController) GetUsers(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	pageNum := ParseInt(query.Get("page"), 1)
	perPage := ParseInt(query.Get("results-per-page"), 10)

	page := NewPage(fake.filteredUsers(query.Get("q")), req.URL.Path, pageNum, perPage)
	response, err := json.Marshal(page)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (fake *CloudController) filteredUsers(query string) *Users {
	switch {
	case strings.Contains(query, "space_guid:"):
		spaceGUID := strings.TrimPrefix(query, "space_guid:")
		space, ok := fake.Spaces.Get(spaceGUID)
		if !ok {
			return NewUsers()
		}

		return space.Developers
	default:
		return fake.Users
	}
}

func ParseInt(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}

	parsedValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		panic(err)
	}

	return int(parsedValue)
}
