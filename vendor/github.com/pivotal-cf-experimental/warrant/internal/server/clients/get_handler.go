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
	token := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
	if ok := h.tokens.Validate(token, []string{"clients"}, []string{"clients.read"}); !ok {
		common.Error(w, http.StatusUnauthorized, "Full authentication is required to access this resource", "unauthorized")
		return
	}

	matches := regexp.MustCompile(`/oauth/clients/(.*)$`).FindStringSubmatch(req.URL.Path)
	id := matches[1]

	client, ok := h.clients.Get(id)
	if !ok {
		common.NotFound(w, fmt.Sprintf("Client %s does not exist", id))
		return
	}

	response, err := json.Marshal(client.ToDocument())
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
