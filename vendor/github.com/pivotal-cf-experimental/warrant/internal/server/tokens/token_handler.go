package tokens

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pivotal-cf-experimental/warrant/internal/server/common"
	"github.com/pivotal-cf-experimental/warrant/internal/server/domain"
)

type urlFinder interface {
	URL() string
}

type tokenHandler struct {
	clients   *domain.Clients
	urlFinder urlFinder
	publicKey string
}

func (h tokenHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// TODO: actually check the basic auth values
	_, _, ok := req.BasicAuth()
	if !ok {
		common.Error(w, http.StatusUnauthorized, "An Authentication object was not found in the SecurityContext", "unauthorized")
		return
	}

	err := req.ParseForm()
	if err != nil {
		panic(err)
	}
	clientID := req.Form.Get("client_id")

	client, ok := h.clients.Get(clientID)
	if !ok {
		panic("client could not be found")
	}

	t := domain.Token{
		ClientID:  clientID,
		Scopes:    client.Scope,
		Audiences: []string{"scim", "password"},
		Issuer:    fmt.Sprintf("%s/oauth/token", h.urlFinder.URL()),
	}

	response, err := json.Marshal(t.ToDocument(h.publicKey))
	if err != nil {
		panic(err)
	}

	w.Write(response)
}
