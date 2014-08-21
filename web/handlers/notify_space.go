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
    notify      Notify
}

func NewNotifySpace(notify Notify, errorWriter ErrorWriterInterface) NotifySpace {
    return NotifySpace{
        errorWriter: errorWriter,
        notify:      notify,
    }
}

func (handler NotifySpace) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    metrics.NewMetric("counter", map[string]interface{}{
        "name": "notifications.web.spaces",
    }).Log()

    transaction := models.NewTransaction()
    err := handler.Execute(w, req, transaction)
    if err != nil {
        handler.errorWriter.Write(w, err)
        return
    }
}

func (handler NotifySpace) Execute(w http.ResponseWriter, req *http.Request, transaction models.TransactionInterface) error {
    transaction.Begin()

    userGUID := postal.SpaceGUID(strings.TrimPrefix(req.URL.Path, "/spaces/"))
    output, err := handler.notify.Execute(transaction, req, userGUID)
    if err != nil {
        transaction.Rollback()
        return err
    }

    transaction.Commit()

    w.WriteHeader(http.StatusOK)
    w.Write(output)

    return nil
}
