package groups

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/pivotal-cf-experimental/warrant/internal/documents"
	"github.com/pivotal-cf-experimental/warrant/internal/server/common"
	"github.com/pivotal-cf-experimental/warrant/internal/server/domain"
)

type addMemberHandler struct {
	groups *domain.Groups
	tokens *domain.Tokens
}

func (h addMemberHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	token := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
	if ok := h.tokens.Validate(token, domain.Token{
		Audiences:   []string{"scim"},
		Authorities: []string{"scim.read"},
	}); !ok {
		common.JSONError(w, http.StatusUnauthorized, "Full authentication is required to access this resource", "unauthorized")
		return
	}

	matches := regexp.MustCompile(`/Groups/(.*)/members$`).FindStringSubmatch(req.URL.Path)
	id := matches[1]

	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	contentType := req.Header.Get("Content-Type")
	if contentType != "application/json" {
		if contentType == "" {
			contentType = http.DetectContentType(requestBody)
		}
		common.JSONError(w, http.StatusBadRequest, fmt.Sprintf("Content type '%s' not supported", contentType), "scim")
		return
	}

	var document documents.CreateMemberRequest
	err = json.Unmarshal(requestBody, &document)
	if err != nil {
		panic(err)
	}

	member, ok := h.groups.AddMember(id, domain.NewMemberFromDocument(document))
	if !ok {
		common.JSONError(w, http.StatusNotFound, fmt.Sprintf("Group %s does not exist", id), "scim_resource_not_found")
		return
	}

	response, err := json.Marshal(member.ToDocument())
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
