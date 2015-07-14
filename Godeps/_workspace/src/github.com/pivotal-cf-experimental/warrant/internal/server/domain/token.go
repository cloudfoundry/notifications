package domain

import (
	"strings"

	"github.com/pivotal-cf-experimental/warrant/internal/documents"
)

type Token struct {
	UserID    string
	ClientID  string
	Scopes    []string
	Audiences []string
}

func newTokenFromClaims(claims map[string]interface{}) Token {
	t := Token{}

	if userID, ok := claims["user_id"].(string); ok {
		t.UserID = userID
	}

	if clientID, ok := claims["client_id"].(string); ok {
		t.ClientID = clientID
	}

	if scopes, ok := claims["scope"].([]interface{}); ok {
		var s []string
		for _, scope := range scopes {
			s = append(s, scope.(string))
		}

		t.Scopes = s
	}

	if audiences, ok := claims["aud"].(string); ok {
		t.Audiences = strings.Split(audiences, " ")
	}

	return t
}

func (t Token) ToDocument() documents.TokenResponse {
	return documents.TokenResponse{
		AccessToken: Tokens{}.Encrypt(t),
		TokenType:   "bearer",
		ExpiresIn:   5000,
		Scope:       strings.Join(t.Scopes, " "),
		JTI:         generateID(),
	}
}

func (t Token) toClaims() map[string]interface{} {
	claims := make(map[string]interface{})

	if len(t.UserID) > 0 {
		claims["user_id"] = t.UserID
	}

	if len(t.ClientID) > 0 {
		claims["client_id"] = t.ClientID
	}

	claims["scope"] = t.Scopes
	claims["aud"] = strings.Join(t.Audiences, " ")

	return claims
}

func (t Token) hasScopes(scopes []string) bool {
	for _, scope := range scopes {
		if !contains(t.Scopes, scope) {
			return false
		}
	}
	return true
}

func (t Token) hasAudiences(audiences []string) bool {
	for _, audience := range audiences {
		if !contains(t.Audiences, audience) {
			return false
		}
	}
	return true
}

func contains(collection []string, item string) bool {
	for _, elem := range collection {
		if elem == item {
			return true
		}
	}

	return false
}
