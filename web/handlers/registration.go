package handlers

import (
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/metrics"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/web/params"
    "github.com/cloudfoundry-incubator/notifications/web/services"
    "github.com/dgrijalva/jwt-go"
)

type Registration struct {
    registrar   services.RegistrarInterface
    errorWriter ErrorWriterInterface
}

func NewRegistration(registrar services.RegistrarInterface, errorWriter ErrorWriterInterface) Registration {
    return Registration{
        registrar:   registrar,
        errorWriter: errorWriter,
    }
}

func (handler Registration) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    metrics.NewMetric("counter", map[string]interface{}{
        "name": "notifications.web.registration",
    }).Log()

    handler.Execute(w, req, models.Database().Connection())
}

func (handler Registration) Execute(w http.ResponseWriter, req *http.Request, connection models.ConnectionInterface) {
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

    client := models.Client{
        ID:          handler.parseClientID(req),
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

func (handler Registration) parseClientID(req *http.Request) string {
    rawToken := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
    clientToken, _ := jwt.Parse(rawToken, func(t *jwt.Token) ([]byte, error) {
        return []byte(config.UAAPublicKey), nil
    })

    return clientToken.Claims["client_id"].(string)
}
