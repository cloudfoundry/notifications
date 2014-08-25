package handlers

import (
    "encoding/json"
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/dgrijalva/jwt-go"
)

type InvalidScopeError string

func (err InvalidScopeError) Error() string {
    return string(err)
}

func NewInvalidScopeError(err string) InvalidScopeError {
    return InvalidScopeError(err)
}

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

func (handler PreferenceFinder) ParseUserID(rawToken string) (string, error) {
    token, err := jwt.Parse(rawToken, func(token *jwt.Token) ([]byte, error) {
        return []byte(config.UAAPublicKey), nil
    })
    if err != nil {
        return "", NewInvalidScopeError(err.Error())
    }

    if scopes, ok := token.Claims["scope"]; ok {
        if handler.HasScope(scopes, "notification_preferences.read") {
            return token.Claims["user_id"].(string), nil
        }
    }

    return "", NewInvalidScopeError("You are not authorized to perform the requested action")
}

func (handler PreferenceFinder) HasScope(elements interface{}, key string) bool {
    for _, elem := range elements.([]interface{}) {
        if elem.(string) == key {
            return true
        }
    }
    return false
}
