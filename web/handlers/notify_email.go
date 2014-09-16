package handlers

import (
    "net/http"

    "github.com/cloudfoundry-incubator/notifications/metrics"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/ryanmoran/stack"
)

type NotifyEmail struct {
    notify NotifyInterface
    recipe postal.MailRecipeInterface
}

func NewNotifyEmail(notify NotifyInterface, errorWriter ErrorWriterInterface, recipe postal.MailRecipeInterface) NotifyEmail {
    return NotifyEmail{
        notify: notify,
        recipe: recipe,
    }
}

func (handler NotifyEmail) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
    connection := models.Database().Connection()
    err := handler.Execute(w, req, connection, context)
    if err != nil {
        panic(err)
    }

    metrics.NewMetric("counter", map[string]interface{}{
        "name": "notifications.web.emails",
    }).Log()
}

func (handler NotifyEmail) Execute(w http.ResponseWriter, req *http.Request, connection models.ConnectionInterface, context stack.Context) error {
    output, err := handler.notify.Execute(connection, req, context, postal.NewEmailID(), handler.recipe)
    if err != nil {
        return err
    }

    w.WriteHeader(http.StatusOK)
    w.Write(output)
    return nil
}
