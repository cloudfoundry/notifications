package handlers

import (
    "net/http"
    "strings"

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
    spaceGUID := postal.UserGUID(strings.TrimPrefix(req.URL.Path, "/users/"))
    output, err := handler.notify.Execute(req, spaceGUID)
    if err != nil {
        handler.errorWriter.Write(w, err)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write(output)
}
