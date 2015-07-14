package tokens

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pivotal-cf-experimental/warrant/internal/server/common"
	"github.com/pivotal-cf-experimental/warrant/internal/server/domain"
)

type authorizeHandler struct {
	tokens *domain.Tokens
	users  *domain.Users
}

func (h authorizeHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Accept") != "application/json" {
		h.redirectToLogin(w)
		return
	}

	requestQuery := req.URL.Query()
	clientID := requestQuery.Get("client_id")
	responseType := requestQuery.Get("response_type")

	if clientID != "cf" {
		h.redirectToLogin(w)
		return
	}

	if responseType != "token" {
		h.redirectToLogin(w)
		return
	}

	req.ParseForm()
	userName := req.Form.Get("username")

	user, ok := h.users.GetByName(userName)
	if !ok {
		common.Error(w, http.StatusNotFound, fmt.Sprintf("User %s does not exist", userName), "scim_resource_not_found")
		return
	}

	if req.Form.Get("source") != "credentials" {
		h.redirectToLogin(w)
		return
	}

	if req.Form.Get("password") != user.Password {
		h.redirectToLogin(w)
		return
	}

	scopes := strings.Join(h.tokens.DefaultScopes, " ")

	t := h.tokens.Encrypt(domain.Token{
		UserID:    user.ID,
		Scopes:    h.tokens.DefaultScopes,
		Audiences: []string{},
	})

	redirectURI := requestQuery.Get("redirect_uri")

	query := url.Values{
		"token_type":   []string{"bearer"},
		"access_token": []string{t},
		"expires_in":   []string{"599"},
		"scope":        []string{scopes},
		"jti":          []string{"ad0efc96-ed29-43ef-be75-85a4b4f105b5"},
	}
	location := fmt.Sprintf("%s#%s", redirectURI, query.Encode())

	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusFound)
}

func (h authorizeHandler) redirectToLogin(w http.ResponseWriter) {
	w.Header().Set("Location", "/login")
	w.WriteHeader(http.StatusFound)
}
