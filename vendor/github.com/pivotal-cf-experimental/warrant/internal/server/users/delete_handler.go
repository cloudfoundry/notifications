package users

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/pivotal-cf-experimental/warrant/internal/server/common"
	"github.com/pivotal-cf-experimental/warrant/internal/server/domain"
)

type deleteHandler struct {
	users  *domain.Users
	tokens *domain.Tokens
}

func (h deleteHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	token := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
	if ok := h.tokens.Validate(token, domain.Token{
		Audiences:   []string{"scim"},
		Authorities: []string{"scim.write"},
	}); !ok {
		common.JSONError(w, http.StatusUnauthorized, "Full authentication is required to access this resource", "unauthorized")
		return
	}

	matches := regexp.MustCompile(`/Users/(.*)$`).FindStringSubmatch(req.URL.Path)
	id := matches[1]

	if ok := h.users.Delete(id); !ok {
		common.JSONError(w, http.StatusNotFound, "User non-existant-user-guid does not exist", "scim_resource_not_found")
		return
	}

	w.WriteHeader(http.StatusOK)
}
