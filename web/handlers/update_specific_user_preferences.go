package handlers

import (
	"net/http"
	"regexp"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/valiant"
	"github.com/cloudfoundry-incubator/notifications/web/params"
	"github.com/ryanmoran/stack"
)

type UpdateSpecificUserPreferences struct {
	preferenceUpdater services.PreferenceUpdaterInterface
	errorWriter       ErrorWriterInterface
}

func NewUpdateSpecificUserPreferences(preferenceUpdater services.PreferenceUpdaterInterface, errorWriter ErrorWriterInterface) UpdateSpecificUserPreferences {
	return UpdateSpecificUserPreferences{
		preferenceUpdater: preferenceUpdater,
		errorWriter:       errorWriter,
	}
}

func (handler UpdateSpecificUserPreferences) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	database := context.Get("database").(models.DatabaseInterface)
	connection := database.Connection()

	userGUID := handler.parseGUID(req.URL.Path)

	builder := services.NewPreferencesBuilder()
	validator := valiant.NewValidator(req.Body)
	err := validator.Validate(&builder)
	if err != nil {
		handler.errorWriter.Write(w, params.ValidationError([]string{err.Error()}))
		return
	}

	preferences, err := builder.ToPreferences()
	if err != nil {
		handler.errorWriter.Write(w, err)
		return
	}

	transaction := connection.Transaction()
	transaction.Begin()
	err = handler.preferenceUpdater.Execute(transaction, preferences, builder.GlobalUnsubscribe, userGUID)
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

func (handler UpdateSpecificUserPreferences) parseGUID(path string) string {
	return regexp.MustCompile(".*/user_preferences/(.*)").FindStringSubmatch(path)[1]
}
