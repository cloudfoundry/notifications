package clients

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/pivotal-cf-experimental/warrant/internal/documents"
	"github.com/pivotal-cf-experimental/warrant/internal/server/common"
	"github.com/pivotal-cf-experimental/warrant/internal/server/domain"
)

type createHandler struct {
	clients *domain.Clients
	tokens  *domain.Tokens
}

func (h createHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
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
		common.JSONError(w, http.StatusUnauthorized, "Full authentication is required to access this resource", "unauthorized")
		return
	}

	var document documents.CreateUpdateClientRequest
	err := json.NewDecoder(req.Body).Decode(&document)
	if err != nil {
		common.HTMLError(w, http.StatusBadRequest, "The request sent by the client was syntactically incorrect.", "")
		return
	}

	client := domain.NewClientFromDocument(document)
	if err := client.Validate(); err != nil {
		common.JSONError(w, http.StatusBadRequest, err.Error(), "invalid_client")
		return
	}

	h.clients.Add(client)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(client.ToDocument())
}
