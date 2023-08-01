package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ryanmoran/stack"
)

type validator interface {
	Parse(string) (*jwt.Token, error)
}

type Authenticator struct {
	Scopes    []string
	Validator validator
}

func NewAuthenticator(validator validator, scopes ...string) Authenticator {
	return Authenticator{
		Scopes:    scopes,
		Validator: validator,
	}
}

func (ware Authenticator) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) bool {
	rawToken := ware.getToken(req)

	if rawToken == "" {
		return ware.Error(w, http.StatusUnauthorized, "Authorization header is invalid: missing")
	}

	token, err := ware.Validator.Parse(rawToken)

	if err != nil {
		return ware.Error(w, http.StatusUnauthorized, "Authorization header is invalid: "+err.Error())
	}

	claims := token.Claims.(jwt.MapClaims)

	if !ware.containsATokenScope(w, token) {
		return false
	}

	context.Set("token", token)
	context.Set("client_id", claims["client_id"])

	return true
}

func (ware Authenticator) containsATokenScope(w http.ResponseWriter, token *jwt.Token) bool {
	claims := token.Claims.(jwt.MapClaims)
	if tokenScopes, ok := claims["scope"]; ok {
		for _, wareScope := range ware.Scopes {
			if contains(tokenScopes, wareScope) {
				return true
			}
		}
	}

	return ware.Error(w, http.StatusForbidden, "You are not authorized to perform the requested action")
}

func contains(elements interface{}, key string) bool {
	for _, elem := range elements.([]interface{}) {
		if elem.(string) == key {
			return true
		}
	}
	return false
}

func (ware Authenticator) Error(w http.ResponseWriter, code int, message string) bool {
	w.WriteHeader(code)
	w.Write([]byte(`{"errors":["` + message + `"]}`))
	return false
}

func (ware Authenticator) getToken(req *http.Request) string {
	authHeader := req.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		return ""
	}

	if strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return parts[1]
}
