package handlers

import (
    "encoding/json"
    "io"
    "io/ioutil"
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/dgrijalva/jwt-go"
)

type Registration struct {
    clientsRepo models.ClientsRepoInterface
    kindsRepo   models.KindsRepoInterface
}

type RegistrationParams struct {
    SourceDescription string        `json:"source_description"`
    Kinds             []models.Kind `json:"kinds"`
}

func NewRegistration(clientsRepo models.ClientsRepoInterface, kindsRepo models.KindsRepoInterface) Registration {
    return Registration{
        clientsRepo: clientsRepo,
        kindsRepo:   kindsRepo,
    }
}

func (handler Registration) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    params, err := handler.parseParams(req.Body)
    if err != nil {
        panic(err)
    }

    client := models.Client{
        ID:          handler.parseClientID(req),
        Description: params.SourceDescription,
    }

    client, err = handler.clientsRepo.Create(client)
    if err != nil {
        panic(err)
    }

    for _, kind := range params.Kinds {
        kind.ClientID = client.ID
        _, err := handler.kindsRepo.Create(kind)

        if err != nil {
            panic(err)
        }
    }
}

func (handler Registration) parseClientID(req *http.Request) string {
    rawToken := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
    clientToken, _ := jwt.Parse(rawToken, func(t *jwt.Token) ([]byte, error) {
        return []byte(config.UAAPublicKey), nil
    })

    return clientToken.Claims["client_id"].(string)
}

func (handler Registration) parseParams(body io.Reader) (RegistrationParams, error) {
    var params RegistrationParams

    bytes, err := ioutil.ReadAll(body)
    if err != nil {
        return params, err
    }

    err = json.Unmarshal(bytes, &params)
    if err != nil {
        return params, err
    }

    return params, nil
}
