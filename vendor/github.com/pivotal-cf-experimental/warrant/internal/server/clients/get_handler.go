package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/pivotal-cf-experimental/warrant/internal/server/common"
	"github.com/pivotal-cf-experimental/warrant/internal/server/domain"
)

type getHandler struct {
	clients *domain.Clients
	tokens  *domain.Tokens
}

func (h getHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
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

	matches := regexp.MustCompile(`/oauth/clients/(.*)$`).FindStringSubmatch(req.URL.Path)
	id := matches[1]

	client, ok := h.clients.Get(id)
	if !ok {
		common.NotFound(w, fmt.Sprintf("Client %s does not exist", id))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(client.ToDocument())
}
