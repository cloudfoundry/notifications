package handlers

import (
    "encoding/json"
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/dgrijalva/jwt-go"
)

type PreferenceFinder struct {
    Preference  PreferenceInterface
    ErrorWriter ErrorWriterInterface
}

func NewPreferenceFinder(preference PreferenceInterface, errorWriter ErrorWriterInterface) PreferenceFinder {
    return PreferenceFinder{
        Preference:  preference,
        ErrorWriter: errorWriter,
    }
}

func (handler PreferenceFinder) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    userID := handler.ParseUserID(strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer "))

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

func (handler PreferenceFinder) ParseUserID(rawToken string) string {
    token, _ := jwt.Parse(rawToken, func(token *jwt.Token) ([]byte, error) {
        return []byte(config.UAAPublicKey), nil
    })

    return token.Claims["user_id"].(string)
}
