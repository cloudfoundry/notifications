package handlers

import (
    "encoding/json"
    "net/http"

    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/cloudfoundry-incubator/notifications/web/params"
    "github.com/cloudfoundry-incubator/notifications/web/services"
    "github.com/dgrijalva/jwt-go"
    "github.com/ryanmoran/stack"
)

type NotifyInterface interface {
    Execute(models.ConnectionInterface, *http.Request, stack.Context, postal.TypedGUID, postal.MailRecipeInterface) ([]byte, error)
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

func (handler Notify) Execute(connection models.ConnectionInterface, req *http.Request, context stack.Context,
    guid postal.TypedGUID, mailRecipe postal.MailRecipeInterface) ([]byte, error) {
    parameters, err := params.NewNotify(req.Body)
    if err != nil {
        return []byte{}, err
    }

    if guid.IsTypeEmail() {
        if !parameters.ValidateEmailRequest() {
            return []byte{}, params.ValidationError(parameters.Errors)
        }
    } else {
        if !parameters.ValidateGUIDRequest() {
            return []byte{}, params.ValidationError(parameters.Errors)
        }
    }

    token := context.Get("token").(*jwt.Token)
    clientID := token.Claims["client_id"].(string)
    client, kind, err := handler.finder.ClientAndKind(clientID, parameters.KindID)
    if err != nil {
        return []byte{}, err
    }

    err = handler.registrar.Register(connection, client, []models.Kind{kind})
    if err != nil {
        return []byte{}, err
    }

    var responses []postal.Response

    responses, err = mailRecipe.Dispatch(clientID, guid, parameters.ToOptions(client, kind), connection)
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
