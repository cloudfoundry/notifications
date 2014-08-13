package handlers

import (
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/web/handlers/params"
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
    handler.Execute(w, req, models.NewTransaction())
}

func (handler Registration) Execute(w http.ResponseWriter, req *http.Request, transaction models.TransactionInterface) {
    parameters, err := params.NewRegistration(req.Body)
    if err != nil {
        handler.errorWriter.Write(w, err)
        return
    }

    err = parameters.Validate()
    if err != nil {
        handler.errorWriter.Write(w, err)
        return
    }

    transaction.Begin()

    client := models.Client{
        ID:          handler.parseClientID(req),
        Description: parameters.SourceDescription,
    }
    client, err = handler.clientsRepo.Upsert(transaction, client)
    if err != nil {
        transaction.Rollback()
        handler.errorWriter.Write(w, err)
        return
    }

    kindIDs := []string{}
    for _, kind := range parameters.Kinds {
        kindIDs = append(kindIDs, kind.ID)

        kind.ClientID = client.ID
        _, err = handler.kindsRepo.Upsert(transaction, kind)
        if err != nil {
            transaction.Rollback()
            handler.errorWriter.Write(w, err)
            return
        }
    }

    if parameters.IncludesKinds {
        _, err = handler.kindsRepo.Trim(transaction, client.ID, kindIDs)
        if err != nil {
            transaction.Rollback()
            handler.errorWriter.Write(w, err)
            return
        }
    }

    transaction.Commit()
}

func (handler Registration) parseClientID(req *http.Request) string {
    rawToken := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
    clientToken, _ := jwt.Parse(rawToken, func(t *jwt.Token) ([]byte, error) {
        return []byte(config.UAAPublicKey), nil
    })

    return clientToken.Claims["client_id"].(string)
}
