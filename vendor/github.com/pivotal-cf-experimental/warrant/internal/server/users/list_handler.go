package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/pivotal-cf-experimental/warrant/internal/server/common"
	"github.com/pivotal-cf-experimental/warrant/internal/server/domain"
)

type listHandler struct {
	users *domain.Users
}

func (h listHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	query, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		panic(err)
	}

	filter := query.Get("filter")
	matches := regexp.MustCompile(`(.*) (.*) ['"](.*)['"]$`).FindStringSubmatch(filter)
	parameter := matches[1]
	operator := matches[2]
	value := matches[3]

	if !validParameter(parameter) {
		common.Error(w, http.StatusBadRequest, fmt.Sprintf("Invalid filter expression: [%s]", filter), "scim")
		return
	}

	if !validOperator(operator) {
		common.Error(w, http.StatusBadRequest, fmt.Sprintf("Invalid filter expression: [%s]", filter), "scim")
		return
	}

	list := domain.UsersList{}

	user, found := h.users.Get(value)
	if found {
		list = append(list, user)
	}

	response, err := json.Marshal(list.ToDocument())
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func validParameter(parameter string) bool {
	for _, p := range []string{"id"} {
		if strings.ToLower(parameter) == p {
			return true
		}
	}

	return false
}

func validOperator(operator string) bool {
	for _, o := range []string{"eq"} {
		if strings.ToLower(operator) == o {
			return true
		}
	}

	return false
}
