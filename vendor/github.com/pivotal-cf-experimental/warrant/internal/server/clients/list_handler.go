package clients

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
	clients *domain.Clients
	tokens  *domain.Tokens
}

// TODO: check for client.admin scope
func (h listHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	token := req.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	token = strings.TrimPrefix(token, "bearer ")

	if len(token) == 0 {
		common.JSONError(w, http.StatusUnauthorized, "Full authentication is required to access this resource", "unauthorized")
		return
	}
	if ok := h.tokens.Validate(token, domain.Token{
		Authorities: []string{"clients.read"},
		Audiences:   []string{"clients"},
	}); !ok {
		common.JSONError(w, http.StatusUnauthorized, "Full authentication is required to access this resource", "unauthorized")
		return
	}

	query, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		panic(err)
	}

	list := domain.ClientsList{}

	filter := query.Get("filter")
	if filter != "" {
		matches := regexp.MustCompile(`(.*) (.*) ['"](.*)['"]$`).FindStringSubmatch(filter)
		parameter := matches[1]
		operator := matches[2]
		value := matches[3]

		if !validParameter(parameter) {
			common.JSONError(w, http.StatusBadRequest, fmt.Sprintf("Invalid filter expression: [%s]", filter), "scim")
			return
		}

		if !validOperator(operator) {
			common.JSONError(w, http.StatusBadRequest, fmt.Sprintf("Invalid filter expression: [%s]", filter), "scim")
			return
		}

		client, found := h.clients.Get(value)
		if found {
			list = append(list, client)
		}
	} else {
		list = append(list, h.clients.All()...)
	}

	switch by := query.Get("sortBy"); by {
	case "name":
		sort.Sort(domain.ByName(list))
	default:
		sort.Sort(domain.ByID(list))
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
