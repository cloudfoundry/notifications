package groups

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/pivotal-cf-experimental/warrant/internal/documents"
	"github.com/pivotal-cf-experimental/warrant/internal/server/common"
	"github.com/pivotal-cf-experimental/warrant/internal/server/domain"
)

type updateHandler struct {
	groups *domain.Groups
	tokens *domain.Tokens
}

func (h updateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	token := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
	if ok := h.tokens.Validate(token, domain.Token{
		Audiences:   []string{"scim"},
		Authorities: []string{"scim.write"},
	}); !ok {
		common.JSONError(w, http.StatusUnauthorized, "Full authentication is required to access this resource", "unauthorized")
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
		common.JSONError(w, http.StatusBadRequest, fmt.Sprintf("Content type '%s' not supported", contentType), "scim")
		return
	}

	var document documents.CreateUpdateGroupRequest
	err = json.Unmarshal(requestBody, &document)
	if err != nil {
		panic(err)
	}

	group := domain.NewGroupFromUpdateDocument(document)

	matches := regexp.MustCompile(`/Groups/(.*)$`).FindStringSubmatch(req.URL.Path)
	id := matches[1]

	existingGroup, ok := h.groups.Get(id)
	if !ok {
		common.JSONError(w, http.StatusNotFound, fmt.Sprintf("Group %s does not exist", group.ID), "scim_resource_not_found")
		return
	}

	version, err := strconv.ParseInt(req.Header.Get("If-Match"), 10, 64)
	if err != nil || existingGroup.Version != int(version) {
		common.JSONError(w, http.StatusBadRequest, "Missing If-Match for PUT", "scim")
		return
	}

	h.groups.Update(group)

	response, err := json.Marshal(group.ToDocument())
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
