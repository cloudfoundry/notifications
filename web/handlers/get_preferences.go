package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/metrics"
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
	metrics.NewMetric("counter", map[string]interface{}{
		"name": "notifications.web.preferences.get",
	}).Log()

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

	result, err := json.Marshal(parsed)
	if err != nil {
		panic(err)
	}

	w.Write(result)
}
