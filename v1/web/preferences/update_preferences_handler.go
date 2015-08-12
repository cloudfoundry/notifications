package preferences

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/cloudfoundry-incubator/notifications/valiant"
	"github.com/cloudfoundry-incubator/notifications/web/webutil"
	"github.com/dgrijalva/jwt-go"
	"github.com/ryanmoran/stack"
)

type UpdatePreferencesHandler struct {
	preferenceUpdater services.PreferenceUpdaterInterface
	errorWriter       errorWriter
}

func NewUpdatePreferencesHandler(preferenceUpdater services.PreferenceUpdaterInterface, errWriter errorWriter) UpdatePreferencesHandler {
	return UpdatePreferencesHandler{
		preferenceUpdater: preferenceUpdater,
		errorWriter:       errWriter,
	}
}

func (h UpdatePreferencesHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	database := context.Get("database").(db.DatabaseInterface)
	connection := database.Connection()

	token := context.Get("token").(*jwt.Token)

	if _, ok := token.Claims["user_id"]; !ok {
		h.errorWriter.Write(w, webutil.MissingUserTokenError("Missing user_id from token claims."))
		return
	}

	userID := token.Claims["user_id"].(string)

	builder := services.NewPreferencesBuilder()
	validator := valiant.NewValidator(req.Body)
	err := validator.Validate(&builder)
	if err != nil {
		h.errorWriter.Write(w, webutil.ValidationError([]string{err.Error()}))
		return
	}

	preferences, err := builder.ToPreferences()
	if err != nil {
		h.errorWriter.Write(w, webutil.ValidationError([]string{err.Error()}))
		return
	}

	transaction := connection.Transaction()
	transaction.Begin()
	err = h.preferenceUpdater.Execute(transaction, preferences, builder.GlobalUnsubscribe, userID)
	if err != nil {
		transaction.Rollback()

		switch err.(type) {
		case services.MissingKindOrClientError, services.CriticalKindError:
			h.errorWriter.Write(w, webutil.ValidationError([]string{err.Error()}))
		default:
			h.errorWriter.Write(w, err)
		}
		return
	}

	err = transaction.Commit()
	if err != nil {
		h.errorWriter.Write(w, models.NewTransactionCommitError(err.Error()))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
