package handlers

import (
    "encoding/json"
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/web/services"
    "github.com/dgrijalva/jwt-go"
)

type GetPreferences struct {
    Preference  services.PreferenceInterface
    ErrorWriter ErrorWriterInterface
}

func NewGetPreferences(preference services.PreferenceInterface, errorWriter ErrorWriterInterface) GetPreferences {
    return GetPreferences{
        Preference:  preference,
        ErrorWriter: errorWriter,
    }
}

func (handler GetPreferences) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    userID, err := handler.ParseUserID(strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer "))
    if err != nil {
        errorWriter := NewErrorWriter()
        errorWriter.Write(w, err)
        return
    }

    parsed, err := handler.Preference.Execute(userID)
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
