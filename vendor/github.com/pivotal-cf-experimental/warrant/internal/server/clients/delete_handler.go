package clients

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/pivotal-cf-experimental/warrant/internal/server/common"
	"github.com/pivotal-cf-experimental/warrant/internal/server/domain"
)

type deleteHandler struct {
	clients *domain.Clients
	tokens  *domain.Tokens
}

func (h deleteHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	token := req.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	token = strings.TrimPrefix(token, "bearer ")

	if len(token) == 0 {
		common.JSONError(w, http.StatusUnauthorized, "Full authentication is required to access this resource", "unauthorized")
		return
	}

	if ok := h.tokens.Validate(token, domain.Token{
		Authorities: []string{"clients.write"},
		Audiences:   []string{"clients"},
	}); !ok {
		common.JSONError(w, http.StatusForbidden, "Invalid token does not contain resource id (clients)", "access_denied")
		return
	}

	matches := regexp.MustCompile(`/oauth/clients/(.*)$`).FindStringSubmatch(req.URL.Path)
	id := matches[1]

	h.clients.Delete(id) // TODO: should return a 404 if the client does not exist
}
