package groups

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/pivotal-cf-experimental/warrant/internal/server/common"
	"github.com/pivotal-cf-experimental/warrant/internal/server/domain"
)

type checkMembershipHandler struct {
	groups *domain.Groups
	tokens *domain.Tokens
}

func (h checkMembershipHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	token := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
	if ok := h.tokens.Validate(token, domain.Token{
		Audiences:   []string{"scim"},
		Authorities: []string{"scim.read"},
	}); !ok {
		common.JSONError(w, http.StatusUnauthorized, "Full authentication is required to access this resource", "unauthorized")
		return
	}

	matches := regexp.MustCompile(`/Groups/(.*)/members/(.*)$`).FindStringSubmatch(req.URL.Path)
	groupID := matches[1]
	memberID := matches[2]

	member, ok := h.groups.CheckMembership(groupID, memberID)
	if !ok {
		common.NotFound(w, fmt.Sprintf("Group %s does not exist or entity %s is not a member", groupID, memberID))
		return
	}

	response, err := json.Marshal(member.ToDocument())
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
