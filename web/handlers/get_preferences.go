package handlers

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/web/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/ryanmoran/stack"
)

type GetPreferences struct {
	PreferencesFinder services.PreferencesFinderInterface
	errorWriter       ErrorWriterInterface
}

func NewGetPreferences(preferencesFinder services.PreferencesFinderInterface, errorWriter ErrorWriterInterface) GetPreferences {
	return GetPreferences{
		PreferencesFinder: preferencesFinder,
		errorWriter:       errorWriter,
	}
}

type MissingUserTokenError string

func (e MissingUserTokenError) Error() string {
	return string(e)
}

func (handler GetPreferences) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	token := context.Get("token").(*jwt.Token)

	if _, ok := token.Claims["user_id"]; !ok {
		handler.errorWriter.Write(w, MissingUserTokenError("Missing user_id from token claims."))
		return
	}

	userID := token.Claims["user_id"].(string)

	parsed, err := handler.PreferencesFinder.Find(userID)
	if err != nil {
		handler.errorWriter.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, parsed)
}
