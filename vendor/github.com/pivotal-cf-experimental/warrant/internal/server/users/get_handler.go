package users

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
	users  *domain.Users
	tokens *domain.Tokens
}

func (h getHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	token := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
	if ok := h.tokens.Validate(token, []string{"scim"}, []string{"scim.read"}); !ok {
		common.Error(w, http.StatusUnauthorized, "Full authentication is required to access this resource", "unauthorized")
		return
	}

	matches := regexp.MustCompile(`/Users/(.*)$`).FindStringSubmatch(req.URL.Path)
	id := matches[1]

	user, ok := h.users.Get(id)
	if !ok {
		common.NotFound(w, fmt.Sprintf("User %s does not exist", id))
		return
	}

	response, err := json.Marshal(user.ToDocument())
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
