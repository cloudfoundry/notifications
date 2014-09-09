package handlers

import (
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/metrics"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
)

type NotifySpace struct {
    errorWriter ErrorWriterInterface
    notify      NotifyInterface
    courier     postal.CourierInterface
}

func NewNotifySpace(notify NotifyInterface, errorWriter ErrorWriterInterface, courier postal.CourierInterface) NotifySpace {
    return NotifySpace{
        errorWriter: errorWriter,
        notify:      notify,
        courier:     courier,
    }
}

func (handler NotifySpace) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    metrics.NewMetric("counter", map[string]interface{}{
        "name": "notifications.web.spaces",
    }).Log()

    connection := models.Database().Connection()
    err := handler.Execute(w, req, connection)
    if err != nil {
        handler.errorWriter.Write(w, err)
        return
    }
}


func (handler NotifySpace) Execute(w http.ResponseWriter, req *http.Request, connection models.ConnectionInterface) error {
    spaceGUID := postal.SpaceGUID(strings.TrimPrefix(req.URL.Path, "/spaces/"))

    output, err := handler.notify.Execute(connection, req, spaceGUID, NewUAARecipe(handler.courier))
    if err != nil {
        return err
    }

    w.WriteHeader(http.StatusOK)
    w.Write(output)

    return nil
}
