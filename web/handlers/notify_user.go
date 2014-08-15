package handlers

import (
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
)

type NotifyUser struct {
    errorWriter ErrorWriterInterface
    notify      Notify
}

func NewNotifyUser(notify Notify, errorWriter ErrorWriterInterface) NotifyUser {
    return NotifyUser{
        errorWriter: errorWriter,
        notify:      notify,
    }
}

func (handler NotifyUser) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    transaction := models.NewTransaction()
    err := handler.Execute(w, req, transaction)
    if err != nil {
        handler.errorWriter.Write(w, err)
        return
    }
}

func (handler NotifyUser) Execute(w http.ResponseWriter, req *http.Request, transaction models.TransactionInterface) error {
    transaction.Begin()

    spaceGUID := postal.UserGUID(strings.TrimPrefix(req.URL.Path, "/users/"))

    output, err := handler.notify.Execute(transaction, req, spaceGUID)
    if err != nil {
        transaction.Rollback()
        return err
    }

    transaction.Commit()

    w.WriteHeader(http.StatusOK)
    w.Write(output)

    return nil
}
