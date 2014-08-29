package handlers

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/web/params"
    "github.com/cloudfoundry-incubator/notifications/web/services"
    "github.com/dgrijalva/jwt-go"
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

func (handler UpdatePreferences) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    userID, err := handler.ParseUserID(strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer "))
    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
        panic(err)
    }

    preferences, err := handler.ParsePreferences(body)
    if err != nil {
        handler.errorWriter.Write(w, params.ParseError{})
        return
    }

    transaction := models.NewTransaction()
    transaction.Begin()
    err = handler.preferenceUpdater.Execute(transaction, preferences, userID)
    if err != nil {
        transaction.Rollback()
        handler.errorWriter.Write(w, err)
        return
    }

    transaction.Commit()
    w.WriteHeader(http.StatusOK)
}

func (handler UpdatePreferences) ParsePreferences(body []byte) ([]models.Preference, error) {
    builder := services.NewPreferencesBuilder()
    err := json.Unmarshal(body, &builder)
    if err != nil {
        return []models.Preference{}, err
    }
    return builder.ToPreferences(), nil
}

func (handler UpdatePreferences) ParseUserID(rawToken string) (string, error) {
    token, err := jwt.Parse(rawToken, func(token *jwt.Token) ([]byte, error) {
        return []byte(config.UAAPublicKey), nil
    })
    if err != nil {
        return "", nil
    }

    return token.Claims["user_id"].(string), nil

}
