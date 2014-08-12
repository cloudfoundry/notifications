package handlers

import (
    "encoding/json"
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/dgrijalva/jwt-go"
)

type NotifySpace struct {
    courier     postal.CourierInterface
    errorWriter ErrorWriterInterface
    clientsRepo models.ClientsRepoInterface
    kindsRepo   models.KindsRepoInterface
}

func NewNotifySpace(courier postal.CourierInterface, errorWriter ErrorWriterInterface,
    clientsRepo models.ClientsRepoInterface, kindsRepo models.KindsRepoInterface) NotifySpace {

    return NotifySpace{
        courier:     courier,
        errorWriter: errorWriter,
        clientsRepo: clientsRepo,
        kindsRepo:   kindsRepo,
    }
}

func (handler NotifySpace) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    params, err := NewNotifyParams(req.Body)
    if err != nil {
        handler.errorWriter.Write(w, err)
        return
    }

    if !params.Validate() {
        handler.errorWriter.Write(w, ParamsValidationError(params.Errors))
        return
    }

    spaceGUID := strings.TrimPrefix(req.URL.Path, "/spaces/")
    rawToken := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")

    token, err := jwt.Parse(rawToken, func(token *jwt.Token) ([]byte, error) {
        return []byte(config.UAAPublicKey), nil
    })

    client, err := handler.FindClient(token.Claims["client_id"].(string))
    if err != nil {
        handler.errorWriter.Write(w, err)
        return
    }

    kind, err := handler.FindKind(params.KindID)
    if err != nil {
        handler.errorWriter.Write(w, err)
        return
    }

    responses, err := handler.courier.Dispatch(rawToken, postal.SpaceGUID(spaceGUID), params.ToOptions(client, kind))
    if err != nil {
        handler.errorWriter.Write(w, err)
        return
    }

    output, err := json.Marshal(responses)
    if err != nil {
        panic(err)
    }

    w.WriteHeader(http.StatusOK)
    w.Write(output)
}

func (handler NotifySpace) FindClient(clientID string) (models.Client, error) {
    client, err := handler.clientsRepo.Find(models.Database().Connection, clientID)
    if err != nil {
        if _, ok := err.(models.ErrRecordNotFound); ok {
            return models.Client{}, nil
        } else {
            return models.Client{}, err
        }
    }
    return client, nil
}

func (handler NotifySpace) FindKind(kindID string) (models.Kind, error) {
    kind, err := handler.kindsRepo.Find(models.Database().Connection, kindID)
    if err != nil {
        if _, ok := err.(models.ErrRecordNotFound); ok {
            return models.Kind{}, nil
        } else {
            return models.Kind{}, err
        }
    }
    return kind, nil
}
