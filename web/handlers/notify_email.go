package handlers

import (
    "net/http"

    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
)

type NotifyEmail struct {
    notify NotifyInterface
    mailer postal.MailerInterface
}

func NewNotifyEmail(notify NotifyInterface, errorWriter ErrorWriterInterface, mailer postal.MailerInterface) NotifyEmail {
    return NotifyEmail{
        notify: notify,
        mailer: mailer,
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
    templateLoader := postal.NewTemplateLoader(postal.NewFileSystem())
    output, err := handler.notify.Execute(connection, req, postal.NewEmailID(), postal.NewEmailRecipe(handler.mailer, templateLoader))
    if err != nil {
        return err
    }

    w.WriteHeader(http.StatusOK)
    w.Write(output)
    return nil
}
