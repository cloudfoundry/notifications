package handlers

import (
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/dgrijalva/jwt-go"
)

type Registration struct {
    clientsRepo models.ClientsRepoInterface
    kindsRepo   models.KindsRepoInterface
    errorWriter ErrorWriterInterface
}

func NewRegistration(clientsRepo models.ClientsRepoInterface, kindsRepo models.KindsRepoInterface, errorWriter ErrorWriterInterface) Registration {
    return Registration{
        clientsRepo: clientsRepo,
        kindsRepo:   kindsRepo,
        errorWriter: errorWriter,
    }
}

func (handler Registration) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    params, err := NewRegistrationParams(req.Body)
    if err != nil {
        handler.errorWriter.Write(w, err)
        return
    }

    err = params.Validate()
    if err != nil {
        handler.errorWriter.Write(w, err)
        return
    }

    client := models.Client{
        ID:          handler.parseClientID(req),
        Description: params.SourceDescription,
    }
    client, err = handler.clientsRepo.Upsert(client)
    if err != nil {
        handler.errorWriter.Write(w, err)
        return
    }

    kindIDs := []string{}
    for _, kind := range params.Kinds {
        kindIDs = append(kindIDs, kind.ID)

        kind.ClientID = client.ID
        _, err = handler.kindsRepo.Upsert(kind)
        if err != nil {
            handler.errorWriter.Write(w, err)
            return
        }
    }

    if params.IncludesKinds {
        _, err = handler.kindsRepo.Trim(client.ID, kindIDs)
        if err != nil {
            handler.errorWriter.Write(w, err)
            return
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
