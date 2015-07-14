package groups

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pivotal-cf-experimental/warrant/internal/documents"
	"github.com/pivotal-cf-experimental/warrant/internal/server/common"
	"github.com/pivotal-cf-experimental/warrant/internal/server/domain"
)

type createHandler struct {
	groups *domain.Groups
	tokens *domain.Tokens
}

func (h createHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	token := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
	if ok := h.tokens.Validate(token, []string{"scim"}, []string{"scim.write"}); !ok {
		common.Error(w, http.StatusUnauthorized, "Full authentication is required to access this resource", "unauthorized")
		return
	}

	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	var document documents.CreateGroupRequest
	err = json.Unmarshal(requestBody, &document)
	if err != nil {
		panic(err)
	}

	if _, ok := h.groups.GetByName(document.DisplayName); ok {
		common.Error(w, http.StatusConflict, fmt.Sprintf("A group with displayName: %s already exists.", document.DisplayName), "scim_resource_already_exists")
		return
	}

	group := domain.NewGroupFromCreateDocument(document)
	h.groups.Add(group)

	response, err := json.Marshal(group.ToDocument())
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
