package handlers

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/valiant"
	"github.com/cloudfoundry-incubator/notifications/web/params"
	"github.com/cloudfoundry-incubator/notifications/web/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/ryanmoran/stack"
)

type UpdatePreferences struct {
	preferenceUpdater services.PreferenceUpdaterInterface
	errorWriter       ErrorWriterInterface
	database          models.DatabaseInterface
}

func NewUpdatePreferences(preferenceUpdater services.PreferenceUpdaterInterface, errorWriter ErrorWriterInterface, database models.DatabaseInterface) UpdatePreferences {
	return UpdatePreferences{
		preferenceUpdater: preferenceUpdater,
		errorWriter:       errorWriter,
		database:          database,
	}
}

func (handler UpdatePreferences) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	metrics.NewMetric("counter", map[string]interface{}{
		"name": "notifications.web.preferences.update",
	}).Log()

	connection := handler.database.Connection()
	handler.Execute(w, req, connection, context)
}

func (handler UpdatePreferences) Execute(w http.ResponseWriter, req *http.Request, connection models.ConnectionInterface, context stack.Context) {
	token := context.Get("token").(*jwt.Token)
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
