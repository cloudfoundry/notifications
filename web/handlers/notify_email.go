package handlers

import (
    "net/http"

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
    err := handler.Execute(w, req, connection)
    if err != nil {
        panic(err)
    }
}

func (handler NotifyEmail) Execute(w http.ResponseWriter, req *http.Request, connection models.ConnectionInterface) error {
    output, err := handler.notify.Execute(connection, req, postal.NewEmailID(), handler.recipe)
    if err != nil {
        return err
    }

    w.WriteHeader(http.StatusOK)
    w.Write(output)
    return nil
}
