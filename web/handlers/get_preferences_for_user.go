package handlers

import (
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/web/services"
	"github.com/ryanmoran/stack"
)

type GetPreferencesForUser struct {
	PreferencesFinder services.PreferencesFinderInterface
	ErrorWriter       ErrorWriterInterface
}

func NewGetPreferencesForUser(preferencesFinder services.PreferencesFinderInterface, errorWriter ErrorWriterInterface) GetPreferencesForUser {
	return GetPreferencesForUser{
		PreferencesFinder: preferencesFinder,
		ErrorWriter:       errorWriter,
	}
}

func (handler GetPreferencesForUser) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	userGUID := handler.parseGUID(req.URL.Path)

	parsed, err := handler.PreferencesFinder.Find(context.Get("database").(models.DatabaseInterface), userGUID)
	if err != nil {
		handler.ErrorWriter.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, parsed)
}

func (handler GetPreferencesForUser) parseGUID(path string) string {
	return strings.TrimPrefix(path, "/user_preferences/")
}
