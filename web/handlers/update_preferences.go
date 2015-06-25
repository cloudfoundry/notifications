package handlers

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/valiant"
	"github.com/cloudfoundry-incubator/notifications/web/params"
	"github.com/dgrijalva/jwt-go"
	"github.com/ryanmoran/stack"
)

type UpdatePreferences struct {
	preferenceUpdater services.PreferenceUpdaterInterface
	errorWriter       ErrorWriterInterface
}

func NewUpdatePreferences(preferenceUpdater services.PreferenceUpdaterInterface, errorWriter ErrorWriterInterface) UpdatePreferences {
	return UpdatePreferences{
		preferenceUpdater: preferenceUpdater,
		errorWriter:       errorWriter,
	}
}

func (handler UpdatePreferences) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	database := context.Get("database").(models.DatabaseInterface)
	connection := database.Connection()

	token := context.Get("token").(*jwt.Token)

	if _, ok := token.Claims["user_id"]; !ok {
		handler.errorWriter.Write(w, MissingUserTokenError("Missing user_id from token claims."))
		return
	}

	userID := token.Claims["user_id"].(string)

	builder := services.NewPreferencesBuilder()
	validator := valiant.NewValidator(req.Body)
	err := validator.Validate(&builder)
	if err != nil {
		handler.errorWriter.Write(w, params.ValidationError([]string{err.Error()}))
		return
	}

	preferences, err := builder.ToPreferences()
	if err != nil {
		handler.errorWriter.Write(w, params.ValidationError([]string{err.Error()}))
		return
	}

	transaction := connection.Transaction()
	transaction.Begin()
	err = handler.preferenceUpdater.Execute(transaction, preferences, builder.GlobalUnsubscribe, userID)
	if err != nil {
		transaction.Rollback()

		switch err.(type) {
		case services.MissingKindOrClientError, services.CriticalKindError:
			handler.errorWriter.Write(w, params.ValidationError([]string{err.Error()}))
		default:
			handler.errorWriter.Write(w, err)
		}
		return
	}

	err = transaction.Commit()
	if err != nil {
		handler.errorWriter.Write(w, models.NewTransactionCommitError(err.Error()))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
