package preferences

import (
	"net/http"
	"strings"

	"github.com/ryanmoran/stack"
)

type GetUserPreferencesHandler struct {
	preferences preferencesFinder
	errorWriter errorWriter
}

func NewGetUserPreferencesHandler(preferences preferencesFinder, errWriter errorWriter) GetUserPreferencesHandler {
	return GetUserPreferencesHandler{
		preferences: preferences,
		errorWriter: errWriter,
	}
}

func (h GetUserPreferencesHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	userGUID := strings.TrimPrefix(req.URL.Path, "/user_preferences/")

	parsed, err := h.preferences.Find(context.Get("database").(DatabaseInterface), userGUID)
	if err != nil {
		h.errorWriter.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, parsed)
}
