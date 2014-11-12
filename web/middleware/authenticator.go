package middleware

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/ryanmoran/stack"
)

type Authenticator struct {
	Scopes       []string
	UAAPublicKey string
}

func NewAuthenticator(publicKey string, scopes ...string) Authenticator {
	return Authenticator{
		Scopes:       scopes,
		UAAPublicKey: publicKey,
	}
}

func (ware Authenticator) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) bool {
	rawToken := ware.getToken(req)

	if rawToken == "" {
		return ware.Error(w, http.StatusUnauthorized, "Authorization header is invalid: missing")
	}

	token, err := jwt.Parse(rawToken, func(t *jwt.Token) ([]byte, error) {
		return []byte(ware.UAAPublicKey), nil
	})
	if err != nil {
		if strings.Contains(err.Error(), "Token is expired") {
			return ware.Error(w, http.StatusUnauthorized, "Authorization header is invalid: expired")
		}
		return ware.Error(w, http.StatusUnauthorized, "Authorization header is invalid: corrupt")
	}

	if scopes, ok := token.Claims["scope"]; ok {
		for _, scope := range ware.Scopes {
			if !ware.HasScope(scopes, scope) {
				return ware.Error(w, http.StatusForbidden, "You are not authorized to perform the requested action")
			}
		}
	} else {
		return ware.Error(w, http.StatusForbidden, "You are not authorized to perform the requested action")
	}

	context.Set("token", token)

	return true
}

func (ware Authenticator) Error(w http.ResponseWriter, code int, message string) bool {
	w.WriteHeader(code)
	w.Write([]byte(`{"errors":["` + message + `"]}`))
	return false
}

func (ware Authenticator) HasScope(elements interface{}, key string) bool {
	for _, elem := range elements.([]interface{}) {
		if elem.(string) == key {
			return true
		}
	}
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
