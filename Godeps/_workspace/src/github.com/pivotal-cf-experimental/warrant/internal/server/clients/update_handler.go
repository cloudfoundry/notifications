package clients

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/pivotal-cf-experimental/warrant/internal/documents"
	"github.com/pivotal-cf-experimental/warrant/internal/server/common"
	"github.com/pivotal-cf-experimental/warrant/internal/server/domain"
)

type updateHandler struct {
	clients *domain.Clients
	tokens  *domain.Tokens
}

func (h updateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	token := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
	if ok := h.tokens.Validate(token, []string{"clients"}, []string{"clients.write"}); !ok {
		common.Error(w, http.StatusUnauthorized, "Full authentication is required to access this resource", "unauthorized")
		return
	}

	var document documents.CreateUpdateClientRequest
	err := json.NewDecoder(req.Body).Decode(&document)
	if err != nil {
		panic(err)
	}

	client := domain.NewClientFromDocument(document)
	if err := client.Validate(); err != nil {
		common.Error(w, http.StatusBadRequest, err.Error(), "invalid_client")
		return
	}

	h.clients.Add(client)

	response, err := json.Marshal(client.ToDocument())
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
