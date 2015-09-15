package middleware

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/ryanmoran/stack"
)

type authenticator interface {
	ServeHTTP(http.ResponseWriter, *http.Request, stack.Context) bool
}

type UnsubscribesAuthenticator struct {
	UAAPublicKey        string
	ClientAuthenticator authenticator
	UserAuthenticator   authenticator
}

func NewUnsubscribesAuthenticator(publicKey string) UnsubscribesAuthenticator {
	return UnsubscribesAuthenticator{
		UAAPublicKey:        publicKey,
		ClientAuthenticator: NewAuthenticator(publicKey, "notification_preferences.admin"),
		UserAuthenticator:   NewAuthenticator(publicKey, "notification_preferences.write"),
	}
}

func (a UnsubscribesAuthenticator) ServeHTTP(writer http.ResponseWriter, request *http.Request, context stack.Context) bool {
	rawToken := a.getToken(request)

	if rawToken == "" {
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write([]byte(`{"errors": [ "Authorization header is invalid: missing" ]}`))
		return false
	}

	token, err := jwt.Parse(rawToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(a.UAAPublicKey), nil
	})
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write([]byte(`{"errors": [ "Authorization header is invalid: corrupt" ]}`))
		return false
	}

	userID, ok := token.Claims["user_id"]
	if ok {
		route := strings.Split(request.URL.Path, "/")
		routeUser := route[len(route)-1]
		if routeUser != userID {
			writer.WriteHeader(http.StatusForbidden)
			writer.Write([]byte(`{"errors": [ "You are not authorized to perform the requested action" ]}`))
			return false
		}
		return a.UserAuthenticator.ServeHTTP(writer, request, context)
	}

	return a.ClientAuthenticator.ServeHTTP(writer, request, context)
}

func (a UnsubscribesAuthenticator) getToken(req *http.Request) string {
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
