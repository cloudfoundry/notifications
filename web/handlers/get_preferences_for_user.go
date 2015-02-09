package handlers

import (
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/web/services"
	"github.com/ryanmoran/stack"
)

type GetPreferencesForUser struct {
	PreferencesFinder services.PreferencesFinderInterface
	ErrorWriter       ErrorWriterInterface
	UserGUID          string
}

func NewGetPreferencesForUser(preferencesFinder services.PreferencesFinderInterface, errorWriter ErrorWriterInterface) GetPreferencesForUser {
	return GetPreferencesForUser{
		PreferencesFinder: preferencesFinder,
		ErrorWriter:       errorWriter,
		UserGUID:          "",
	}
}

func (handler GetPreferencesForUser) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	metrics.NewMetric("counter", map[string]interface{}{
		"name": "notifications.web.preferences.get",
	}).Log()

	handler.UserGUID = handler.parseGUID(req.URL.Path)

	parsed, err := handler.PreferencesFinder.Find(handler.UserGUID)
	if err != nil {
		handler.ErrorWriter.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, parsed)
}

func (handler GetPreferencesForUser) parseGUID(path string) string {
	return strings.TrimPrefix(path, "/user_preferences/")
}
