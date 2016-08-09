package groups

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/pivotal-cf-experimental/warrant/internal/server/common"
	"github.com/pivotal-cf-experimental/warrant/internal/server/domain"
)

type deleteHandler struct {
	groups *domain.Groups
	tokens *domain.Tokens
}

func (h deleteHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	token := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
	if ok := h.tokens.Validate(token, []string{"scim"}, []string{"scim.write"}); !ok {
		common.Error(w, http.StatusUnauthorized, "Full authentication is required to access this resource", "unauthorized")
		return
	}

	matches := regexp.MustCompile(`/Groups/(.*)$`).FindStringSubmatch(req.URL.Path)
	id := matches[1]

	if ok := h.groups.Delete(id); !ok {
		common.Error(w, http.StatusNotFound, fmt.Sprintf("Group %s does not exist", id), "scim_resource_not_found")
		return
	}

	w.WriteHeader(http.StatusOK)
}
