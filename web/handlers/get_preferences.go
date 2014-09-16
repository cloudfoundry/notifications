package handlers

import (
    "encoding/json"
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/config"
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
    userID, err := handler.ParseUserID(strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer "))
    if err != nil {
        errorWriter := NewErrorWriter()
        errorWriter.Write(w, err)
        return
    }

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
}

func (handler GetPreferences) ParseUserID(rawToken string) (string, error) {
    token, err := jwt.Parse(rawToken, func(token *jwt.Token) ([]byte, error) {
        return []byte(config.UAAPublicKey), nil
    })
    return token.Claims["user_id"].(string), err
}
