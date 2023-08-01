package preferences

import (
	"errors"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/cloudfoundry-incubator/notifications/v1/web/webutil"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ryanmoran/stack"
)

type errorWriter interface {
	Write(writer http.ResponseWriter, err error)
}

type preferencesFinder interface {
	Find(database services.DatabaseInterface, userGUID string) (services.PreferencesBuilder, error)
}

type GetPreferencesHandler struct {
	preferences preferencesFinder
	errorWriter errorWriter
}

func NewGetPreferencesHandler(preferences preferencesFinder, errWriter errorWriter) GetPreferencesHandler {
	return GetPreferencesHandler{
		preferences: preferences,
		errorWriter: errWriter,
	}
}

func (h GetPreferencesHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	token := context.Get("token").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	if _, ok := claims["user_id"]; !ok {
		h.errorWriter.Write(w, webutil.MissingUserTokenError{Err: errors.New("Missing user_id from token claims.")})
		return
	}

	userID := claims["user_id"].(string)

	parsed, err := h.preferences.Find(context.Get("database").(DatabaseInterface), userID)
	if err != nil {
		h.errorWriter.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, parsed)
}
