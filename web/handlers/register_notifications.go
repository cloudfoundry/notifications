package handlers

import (
    "net/http"

    "github.com/cloudfoundry-incubator/notifications/metrics"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/web/params"
    "github.com/cloudfoundry-incubator/notifications/web/services"
    "github.com/dgrijalva/jwt-go"
    "github.com/ryanmoran/stack"
)

type RegisterNotifications struct {
    registrar   services.RegistrarInterface
    errorWriter ErrorWriterInterface
}

func NewRegisterNotifications(registrar services.RegistrarInterface, errorWriter ErrorWriterInterface) RegisterNotifications {
    return RegisterNotifications{
        registrar:   registrar,
        errorWriter: errorWriter,
    }
}

func (handler RegisterNotifications) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
    handler.Execute(w, req, models.Database().Connection(), context)

    metrics.NewMetric("counter", map[string]interface{}{
        "name": "notifications.web.registration",
    }).Log()
}

func (handler RegisterNotifications) Execute(w http.ResponseWriter, req *http.Request, connection models.ConnectionInterface, context stack.Context) {
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

    transaction := connection.Transaction()
    transaction.Begin()

    token := context.Get("token").(*jwt.Token)
    clientID := token.Claims["client_id"].(string)

    client := models.Client{
        ID:          clientID,
        Description: parameters.SourceDescription,
    }

    kinds := []models.Kind{}
    for _, kind := range parameters.Kinds {
        kind.ClientID = client.ID
        kinds = append(kinds, kind)
    }

    err = handler.registrar.Register(transaction, client, kinds)
    if err != nil {
        transaction.Rollback()
        handler.errorWriter.Write(w, err)
        return
    }

    if parameters.IncludesKinds {
        err = handler.registrar.Prune(transaction, client, kinds)
        if err != nil {
            transaction.Rollback()
            handler.errorWriter.Write(w, err)
            return
        }
    }

    transaction.Commit()
}
