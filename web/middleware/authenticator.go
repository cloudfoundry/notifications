package middleware

import (
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/dgrijalva/jwt-go"
)

type Authenticator struct{}

func NewAuthenticator() Authenticator {
    return Authenticator{}
}

func (ware Authenticator) ServeHTTP(w http.ResponseWriter, req *http.Request) bool {
    authHeader := req.Header.Get("Authorization")
    rawToken := strings.TrimPrefix(authHeader, "Bearer ")

    if rawToken == "" {
        return ware.Error(w, http.StatusUnauthorized, "Authorization header is invalid: missing")
    }

    token, err := jwt.Parse(rawToken, func(t *jwt.Token) ([]byte, error) {
        return []byte(config.UAAPublicKey), nil
    })
    if err != nil {
        if strings.Contains(err.Error(), "Token is expired") {
            return ware.Error(w, http.StatusUnauthorized, "Authorization header is invalid: expired")
        }
        return ware.Error(w, http.StatusUnauthorized, "Authorization header is invalid: corrupt")
    }

    if scopes, ok := token.Claims["scope"]; ok {
        if !ware.HasScope(scopes, "notifications.write") {
            return ware.Error(w, http.StatusForbidden, "You are not authorized to perform the requested action")
        }
    } else {
        return ware.Error(w, http.StatusForbidden, "You are not authorized to perform the requested action")
    }

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
