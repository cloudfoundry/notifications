package users

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
	users  *domain.Users
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

	contentType := req.Header.Get("Content-Type")
	if contentType != "application/json" {
		if contentType == "" {
			contentType = http.DetectContentType(requestBody)
		}
		common.Error(w, http.StatusBadRequest, fmt.Sprintf("Content type '%s' not supported", contentType), "scim")
		return
	}

	var document documents.CreateUserRequest
	err = json.Unmarshal(requestBody, &document)
	if err != nil {
		panic(err)
	}

	if _, ok := h.users.GetByName(document.UserName); ok {
		common.Error(w, http.StatusConflict, fmt.Sprintf("Username already in use: %s", document.UserName), "scim_resource_already_exists")
		return
	}

	user := domain.NewUserFromCreateDocument(document)
	if err := user.Validate(); err != nil {
		common.Error(w, http.StatusBadRequest, err.Error(), "invalid_scim_resource")
		return
	}
	h.users.Add(user)

	response, err := json.Marshal(user.ToDocument())
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
