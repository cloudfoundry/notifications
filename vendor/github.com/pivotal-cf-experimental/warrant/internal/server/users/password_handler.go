package users

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/pivotal-cf-experimental/warrant/internal/documents"
	"github.com/pivotal-cf-experimental/warrant/internal/server/common"
	"github.com/pivotal-cf-experimental/warrant/internal/server/domain"
)

type passwordHandler struct {
	users  *domain.Users
	tokens *domain.Tokens
}

func (h passwordHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	token := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
	matches := regexp.MustCompile(`/Users/(.*)/password$`).FindStringSubmatch(req.URL.Path)
	id := matches[1]

	user, ok := h.users.Get(id)
	if !ok {
		common.JSONError(w, http.StatusUnauthorized, "Not authorized", "access_denied")
		return
	}

	var document documents.ChangePasswordRequest
	err := json.NewDecoder(req.Body).Decode(&document)
	if err != nil {
		panic(err)
	}

	if !h.canUpdateUserPassword(id, token, user.Password, document.OldPassword) {
		common.JSONError(w, http.StatusUnauthorized, "Not authorized", "access_denied")
		return
	}

	user.Password = document.Password
	h.users.Update(user)
}

func (h passwordHandler) canUpdateUserPassword(userID, tokenHeader, existingPassword, givenPassword string) bool {
	if h.tokens.Validate(tokenHeader, domain.Token{
		Audiences:   []string{"password"},
		Authorities: []string{"password.write"},
	}) {
		return true
	}

	t, err := h.tokens.Decrypt(tokenHeader)
	if err != nil {
		return false
	}

	if t.UserID == userID && existingPassword == givenPassword {
		return true
	}

	return false
}
