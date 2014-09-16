package handlers

import (
    "encoding/json"
    "net/http"

    "github.com/cloudfoundry-incubator/notifications/metrics"
    "github.com/cloudfoundry-incubator/notifications/web/services"
    "github.com/dgrijalva/jwt-go"
    "github.com/ryanmoran/stack"
)

type GetPreferences struct {
    PreferencesFinder services.PreferencesFinderInterface
    ErrorWriter       ErrorWriterInterface
}

func NewGetPreferences(preferencesFinder services.PreferencesFinderInterface, errorWriter ErrorWriterInterface) GetPreferences {
    return GetPreferences{
        PreferencesFinder: preferencesFinder,
        ErrorWriter:       errorWriter,
    }
}

func (handler GetPreferences) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
    token := context.Get("token").(*jwt.Token)
    userID := token.Claims["user_id"].(string)

    parsed, err := handler.PreferencesFinder.Find(userID)
    if err != nil {
        errorWriter := NewErrorWriter()
        errorWriter.Write(w, err)
        return
    }

    result, err := json.Marshal(parsed)
    if err != nil {
        panic(err)
    }

    w.Write(result)

    metrics.NewMetric("counter", map[string]interface{}{
        "name": "notifications.web.preferences.get",
    }).Log()
}
