package handlers

import (
    "net/http"

    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
)

type NotifyEmail struct {
    notify  NotifyInterface
    courier postal.CourierInterface
}

func NewNotifyEmail(notify NotifyInterface, errorWriter ErrorWriterInterface, courier postal.CourierInterface) NotifyEmail {
    return NotifyEmail{
        notify:  notify,
        courier: courier,
    }
}

func (handler NotifyEmail) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    connection := models.Database().Connection()

    err := handler.Execute(w, req, connection)
    if err != nil {
        panic(err)
    }
}

func (handler NotifyEmail) Execute(w http.ResponseWriter, req *http.Request, connection models.ConnectionInterface) error {
    output, err := handler.notify.Execute(connection, req, postal.NewEmailID(), NewEmailRecipe(handler.courier))
    if err != nil {
        return err
    }

    w.WriteHeader(http.StatusOK)
    w.Write(output)
    return nil
}
