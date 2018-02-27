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
	clients    *domain.Clients
	users      *domain.Users
	urlFinder  urlFinder
	privateKey string
}

func (h tokenHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// TODO: actually check the basic auth values
	clientID, _, ok := req.BasicAuth()
	if !ok {
		common.JSONError(w, http.StatusUnauthorized, "An Authentication object was not found in the SecurityContext", "unauthorized")
		return
	}

	client, ok := h.clients.Get(clientID)
	if !ok {
		common.JSONError(w, http.StatusUnauthorized, fmt.Sprintf("No client with requested id: %s", clientID), "invalid_client")
		return
	}

	err := req.ParseForm()
	if err != nil {
		panic(err)
	}

	var t domain.Token
	if req.Form.Get("grant_type") == "client_credentials" {
		t.ClientID = clientID
		t.Scopes = client.Scope
		t.Authorities = client.Authorities
		t.Audiences = client.ResourceIDs
		t.Issuer = fmt.Sprintf("%s/oauth/token", h.urlFinder.URL())
	} else {
		user, ok := h.users.GetByName(req.Form.Get("username"))
		if !ok {
			common.JSONError(w, http.StatusNotFound, fmt.Sprintf("User %s does not exist", req.Form.Get("username")), "scim_resource_not_found")
			return
		}

		t.ClientID = clientID
		t.Scopes = client.Scope
		t.UserID = user.ID
	}

	response, err := json.Marshal(t.ToDocument(h.privateKey))
	if err != nil {
		panic(err)
	}

	w.Write(response)
}
