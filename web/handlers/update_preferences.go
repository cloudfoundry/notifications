package handlers

import (
    "encoding/json"
    "io/ioutil"
    "net/http"

    "github.com/cloudfoundry-incubator/notifications/metrics"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/web/params"
    "github.com/cloudfoundry-incubator/notifications/web/services"
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
    connection := models.Database().Connection()
    handler.Execute(w, req, connection, context)

    metrics.NewMetric("counter", map[string]interface{}{
        "name": "notifications.web.preferences.update",
    }).Log()
}

func (handler UpdatePreferences) Execute(w http.ResponseWriter, req *http.Request, connection models.ConnectionInterface, context stack.Context) {
    token := context.Get("token").(*jwt.Token)
    userID := token.Claims["user_id"].(string)

    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
        panic(err)
    }

    builder, err := handler.ParsePreferences(body)
    if err != nil {
        handler.errorWriter.Write(w, params.ParseError{})
        return
    }
    preferences, err := builder.ToPreferences()
    if err != nil {
        handler.errorWriter.Write(w, params.ValidationError{})
        return
    }

    transaction := connection.Transaction()
    transaction.Begin()
    err = handler.preferenceUpdater.Execute(transaction, preferences, userID)
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

func (handler UpdatePreferences) ParsePreferences(body []byte) (services.PreferencesBuilder, error) {
    builder := services.NewPreferencesBuilder()
    err := json.Unmarshal(body, &builder)
    if err != nil {
        return builder, err
    }
    return builder, nil
}
