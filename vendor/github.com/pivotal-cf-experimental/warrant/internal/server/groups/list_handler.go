package groups

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"

	"github.com/pivotal-cf-experimental/warrant/internal/server/common"
	"github.com/pivotal-cf-experimental/warrant/internal/server/domain"
)

type listHandler struct {
	groups *domain.Groups
	tokens *domain.Tokens
}

func (h listHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	token := req.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	token = strings.TrimPrefix(token, "bearer ")

	if len(token) == 0 {
		common.JSONError(w, http.StatusUnauthorized, "Full authentication is required to access this resource", "unauthorized")
		return
	}
	if ok := h.tokens.Validate(token, domain.Token{
		Authorities: []string{"scim.read"},
		Audiences:   []string{"scim"},
	}); !ok {
		common.JSONError(w, http.StatusUnauthorized, "Full authentication is required to access this resource", "unauthorized")
		return
	}

	query, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		panic(err)
	}

	list := domain.GroupsList{}

	filter := query.Get("filter")
	if filter != "" {
		matches := regexp.MustCompile(`(.*) (.*) ['"](.*)['"]$`).FindStringSubmatch(filter)
		parameter := matches[1]
		operator := matches[2]
		value := matches[3]

		if !validParameter(parameter) {
			common.JSONError(w, http.StatusBadRequest, fmt.Sprintf("Invalid parameter in filter expression: [%s]", filter), "scim")
			return
		}

		if !validOperator(operator) {
			common.JSONError(w, http.StatusBadRequest, fmt.Sprintf("Invalid operator in filter expression: [%s]", filter), "scim")
			return
		}

		if strings.ToLower(parameter) == "displayname" {
			group, found := h.groups.GetByName(value)
			if found {
				list = append(list, group)
			}
		} else {
			group, found := h.groups.Get(value)
			if found {
				list = append(list, group)
			}
		}
	} else {
		list = h.groups.All()
	}

	by := strings.ToLower(query.Get("sortBy"))
	if by == "displayname" {
		sort.Sort(domain.GroupsByDisplayName(list))
	} else {
		sort.Sort(domain.GroupsByCreatedAt(list))
	}

	response, err := json.Marshal(list.ToDocument())
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func validParameter(parameter string) bool {
	for _, p := range []string{"id", "displayname"} {
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
