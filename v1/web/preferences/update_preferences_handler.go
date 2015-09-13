package preferences

import (
	"errors"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/cloudfoundry-incubator/notifications/v1/web/webutil"
	"github.com/cloudfoundry-incubator/notifications/valiant"
	"github.com/dgrijalva/jwt-go"
	"github.com/ryanmoran/stack"
)

type UpdatePreferencesHandler struct {
	preferences preferenceUpdater
	errorWriter errorWriter
}

func NewUpdatePreferencesHandler(preferences preferenceUpdater, errWriter errorWriter) UpdatePreferencesHandler {
	return UpdatePreferencesHandler{
		preferences: preferences,
		errorWriter: errWriter,
	}
}

func (h UpdatePreferencesHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	database := context.Get("database").(DatabaseInterface)
	connection := database.Connection()

	token := context.Get("token").(*jwt.Token)

	if _, ok := token.Claims["user_id"]; !ok {
		h.errorWriter.Write(w, webutil.MissingUserTokenError{errors.New("Missing user_id from token claims.")})
		return
	}

	userID := token.Claims["user_id"].(string)

	builder := services.NewPreferencesBuilder()
	validator := valiant.NewValidator(req.Body)
	err := validator.Validate(&builder)
	if err != nil {
		h.errorWriter.Write(w, webutil.ValidationError{err})
		return
	}

	preferences, err := builder.ToPreferences()
	if err != nil {
		h.errorWriter.Write(w, webutil.ValidationError{err})
		return
	}

	transaction := connection.Transaction()
	transaction.Begin()
	err = h.preferences.Update(transaction, preferences, builder.GlobalUnsubscribe, userID)
	if err != nil {
		transaction.Rollback()

		switch err.(type) {
		case services.MissingKindOrClientError, services.CriticalKindError:
			h.errorWriter.Write(w, webutil.ValidationError{err})
		default:
			h.errorWriter.Write(w, err)
		}
		return
	}

	err = transaction.Commit()
	if err != nil {
		h.errorWriter.Write(w, models.TransactionCommitError{err})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
