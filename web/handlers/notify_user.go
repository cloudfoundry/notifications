package handlers

import (
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/metrics"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
)

type NotifyUser struct {
    errorWriter ErrorWriterInterface
    notify      NotifyInterface
    courier     postal.CourierInterface
}

func NewNotifyUser(notify NotifyInterface, errorWriter ErrorWriterInterface, courier postal.CourierInterface) NotifyUser {
    return NotifyUser{
        errorWriter: errorWriter,
        notify:      notify,
        courier:     courier,
    }
}

func (handler NotifyUser) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    metrics.NewMetric("counter", map[string]interface{}{
        "name": "notifications.web.users",
    }).Log()

    connection := models.Database().Connection()
    err := handler.Execute(w, req, connection)
    if err != nil {
        handler.errorWriter.Write(w, err)
        return
    }
}

func (handler NotifyUser) Execute(w http.ResponseWriter, req *http.Request, connection models.ConnectionInterface) error {
    userGUID := postal.UserGUID(strings.TrimPrefix(req.URL.Path, "/users/"))

    output, err := handler.notify.Execute(connection, req, userGUID, postal.NewUAARecipe(handler.courier))
    if err != nil {
        return err
    }

    w.WriteHeader(http.StatusOK)
    w.Write(output)

    return nil
}
