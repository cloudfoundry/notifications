package preferences

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/webutil"
	"github.com/dgrijalva/jwt-go"
	"github.com/ryanmoran/stack"
)

type errorWriter interface {
	Write(writer http.ResponseWriter, err error)
}

type GetPreferencesHandler struct {
	preferencesFinder services.PreferencesFinderInterface
	errorWriter       errorWriter
}

func NewGetPreferencesHandler(preferencesFinder services.PreferencesFinderInterface, errWriter errorWriter) GetPreferencesHandler {
	return GetPreferencesHandler{
		preferencesFinder: preferencesFinder,
		errorWriter:       errWriter,
	}
}

func (h GetPreferencesHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	token := context.Get("token").(*jwt.Token)

	if _, ok := token.Claims["user_id"]; !ok {
		h.errorWriter.Write(w, webutil.MissingUserTokenError("Missing user_id from token claims."))
		return
	}

	userID := token.Claims["user_id"].(string)

	parsed, err := h.preferencesFinder.Find(context.Get("database").(models.DatabaseInterface), userID)
	if err != nil {
		h.errorWriter.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, parsed)
}
