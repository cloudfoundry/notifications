package handlers

import (
    "encoding/json"
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/cloudfoundry-incubator/notifications/web/params"
    "github.com/cloudfoundry-incubator/notifications/web/services"
    "github.com/dgrijalva/jwt-go"
)

type NotifyInterface interface {
    Execute(models.ConnectionInterface, *http.Request, postal.TypedGUID, postal.MailRecipeInterface) ([]byte, error)
}

type Notify struct {
    finder    services.NotificationFinderInterface
    registrar services.RegistrarInterface
}

func NewNotify(finder services.NotificationFinderInterface, registrar services.RegistrarInterface) Notify {
    return Notify{
        finder:    finder,
        registrar: registrar,
    }
}

func (handler Notify) Execute(connection models.ConnectionInterface, req *http.Request,
    guid postal.TypedGUID, mailRecipe postal.MailRecipeInterface) ([]byte, error) {
    parameters, err := params.NewNotify(req.Body)
    if err != nil {
        return []byte{}, err
    }

    if !parameters.Validate(guid) {
        return []byte{}, params.ValidationError(parameters.Errors)
    }

    clientID := handler.ParseClientID(strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer "))
    client, kind, err := handler.finder.ClientAndKind(clientID, parameters.KindID)
    if err != nil {
        return []byte{}, err
    }

    err = handler.registrar.Register(connection, client, []models.Kind{kind})
    if err != nil {
        return []byte{}, err
    }

    var responses []postal.Response

    responses, err = mailRecipe.DeliverMail(clientID, guid, parameters.ToOptions(client, kind), connection)
    if err != nil {
        return []byte{}, err
    }

    output, err := json.Marshal(responses)
    if err != nil {
        panic(err)
    }

    output = mailRecipe.Trim(output)

    return output, nil
}

func (handler Notify) ParseClientID(rawToken string) string {
    token, _ := jwt.Parse(rawToken, func(token *jwt.Token) ([]byte, error) {
        return []byte(config.UAAPublicKey), nil
    })
    return token.Claims["client_id"].(string)
}
