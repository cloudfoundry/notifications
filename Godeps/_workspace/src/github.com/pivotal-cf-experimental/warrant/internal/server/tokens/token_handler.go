package tokens

import (
	"encoding/json"
	"net/http"

	"github.com/pivotal-cf-experimental/warrant/internal/server/common"
	"github.com/pivotal-cf-experimental/warrant/internal/server/domain"
)

type tokenHandler struct{}

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

	scopes := []string{"scim.write", "scim.read", "password.write"}
	t := domain.Token{
		ClientID:  clientID,
		Scopes:    scopes,
		Audiences: []string{"scim", "password"},
	}

	response, err := json.Marshal(t.ToDocument())
	if err != nil {
		panic(err)
	}

	w.Write(response)
}
