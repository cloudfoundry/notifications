package handlers

import (
    "net/http"
    "strings"

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
    userGUID := postal.SpaceGUID(strings.TrimPrefix(req.URL.Path, "/spaces/"))
    output, err := handler.notify.Execute(req, userGUID)
    if err != nil {
        handler.errorWriter.Write(w, err)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write(output)
}
