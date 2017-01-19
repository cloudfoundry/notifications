package middleware

import (
	"net/http"
	"strings"

	"github.com/ryanmoran/stack"
)

type authenticator interface {
	ServeHTTP(http.ResponseWriter, *http.Request, stack.Context) bool
}

type UnsubscribesAuthenticator struct {
	ClientAuthenticator authenticator
	UserAuthenticator   authenticator
	Validator           validator
}

func NewUnsubscribesAuthenticator(validator validator) UnsubscribesAuthenticator {
	return UnsubscribesAuthenticator{
		Validator:           validator,
		ClientAuthenticator: NewAuthenticator(validator, "notification_preferences.admin"),
		UserAuthenticator:   NewAuthenticator(validator, "notification_preferences.write"),
	}
}

func (a UnsubscribesAuthenticator) ServeHTTP(writer http.ResponseWriter, request *http.Request, context stack.Context) bool {
	rawToken := a.getToken(request)

	if rawToken == "" {
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write([]byte(`{"errors": [ "Authorization header is invalid: missing" ]}`))
		return false
	}

	token, err := a.Validator.Parse(rawToken)

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
