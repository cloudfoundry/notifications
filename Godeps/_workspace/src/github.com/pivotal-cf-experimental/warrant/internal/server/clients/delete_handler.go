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
	token := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
	if len(token) == 0 {
		common.Error(w, http.StatusUnauthorized, "Full authentication is required to access this resource", "unauthorized")
		return
	}

	if ok := h.tokens.Validate(token, []string{"clients"}, []string{"clients.write"}); !ok {
		common.Error(w, http.StatusForbidden, "Invalid token does not contain resource id (clients)", "access_denied")
		return
	}

	matches := regexp.MustCompile(`/oauth/clients/(.*)$`).FindStringSubmatch(req.URL.Path)
	id := matches[1]

	if ok := h.clients.Delete(id); !ok {
		panic("foo")
	}

	w.WriteHeader(http.StatusOK)
}
