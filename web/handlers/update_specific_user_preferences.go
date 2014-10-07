package handlers

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
    "regexp"

    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/web/params"
    "github.com/cloudfoundry-incubator/notifications/web/services"
    "github.com/ryanmoran/stack"
)

type UpdateSpecificUserPreferences struct {
    preferenceUpdater services.PreferenceUpdaterInterface
    errorWriter       ErrorWriterInterface
    database          models.DatabaseInterface
}

func NewUpdateSpecificUserPreferences(preferenceUpdater services.PreferenceUpdaterInterface, errorWriter ErrorWriterInterface, database models.DatabaseInterface) UpdateSpecificUserPreferences {
    return UpdateSpecificUserPreferences{
        preferenceUpdater: preferenceUpdater,
        errorWriter:       errorWriter,
        database:          database,
    }
}

func (handler UpdateSpecificUserPreferences) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
    connection := handler.database.Connection()
    handler.Execute(w, req, connection, context)
}

func (handler UpdateSpecificUserPreferences) Execute(w http.ResponseWriter, req *http.Request, conn models.ConnectionInterface, context stack.Context) {
    userGUID := handler.parseGUID(req.URL.Path)

    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
        panic(err)
    }

    builder, err := handler.ParsePreferences(body)
    if err != nil {
        panic(err)
    }

    preferences, err := builder.ToPreferences()
    if err != nil {
        panic(err)
    }

    transaction := conn.Transaction()
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
    transaction.Commit()

    w.WriteHeader(http.StatusNoContent)
}

func (handler UpdateSpecificUserPreferences) parseGUID(path string) string {

    regex := regexp.MustCompile(".*/user_preferences/(.*)")

    return regex.FindStringSubmatch(path)[1]
}

func (handler UpdateSpecificUserPreferences) ParsePreferences(body []byte) (services.PreferencesBuilder, error) {
    builder := services.NewPreferencesBuilder()
    err := json.Unmarshal(body, &builder)
    if err != nil {
        return builder, err
    }
    return builder, nil
}
