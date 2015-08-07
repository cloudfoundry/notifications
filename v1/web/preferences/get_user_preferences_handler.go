package preferences

import (
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/ryanmoran/stack"
)

type GetUserPreferencesHandler struct {
	preferencesFinder services.PreferencesFinderInterface
	errorWriter       errorWriter
}

func NewGetUserPreferencesHandler(preferencesFinder services.PreferencesFinderInterface, errWriter errorWriter) GetUserPreferencesHandler {
	return GetUserPreferencesHandler{
		preferencesFinder: preferencesFinder,
		errorWriter:       errWriter,
	}
}

func (h GetUserPreferencesHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	userGUID := strings.TrimPrefix(req.URL.Path, "/user_preferences/")

	parsed, err := h.preferencesFinder.Find(context.Get("database").(models.DatabaseInterface), userGUID)
	if err != nil {
		h.errorWriter.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, parsed)
}
