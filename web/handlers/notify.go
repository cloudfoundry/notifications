package handlers

import (
    "encoding/json"
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/cloudfoundry-incubator/notifications/web/handlers/params"
    "github.com/dgrijalva/jwt-go"
)

type Notify struct {
    courier   postal.CourierInterface
    finder    FinderInterface
    registrar RegistrarInterface
}

func NewNotify(courier postal.CourierInterface, finder FinderInterface, registrar RegistrarInterface) Notify {
    return Notify{
        courier:   courier,
        finder:    finder,
        registrar: registrar,
    }
}

func (handler Notify) Execute(conn models.ConnectionInterface, req *http.Request, guid postal.TypedGUID) ([]byte, error) {
    parameters, err := params.NewNotify(req.Body)
    if err != nil {
        return []byte{}, err
    }

    if !parameters.Validate() {
        return []byte{}, params.ValidationError(parameters.Errors)
    }

    clientID := handler.ParseClientID(strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer "))
    client, kind, err := handler.finder.ClientAndKind(clientID, parameters.KindID)
    if err != nil {
        return []byte{}, err
    }

    err = handler.registrar.Register(conn, client, []models.Kind{kind})
    if err != nil {
        return []byte{}, err
    }

    responses, err := handler.courier.Dispatch(clientID, guid, parameters.ToOptions(client, kind), conn)
    if err != nil {
        return []byte{}, err
    }

    output, err := json.Marshal(responses)
    if err != nil {
        panic(err)
    }

    return output, nil
}

func (handler Notify) ParseClientID(rawToken string) string {
    token, _ := jwt.Parse(rawToken, func(token *jwt.Token) ([]byte, error) {
        return []byte(config.UAAPublicKey), nil
    })
    return token.Claims["client_id"].(string)
}
