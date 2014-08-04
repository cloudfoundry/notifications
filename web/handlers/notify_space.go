package handlers

import (
    "encoding/json"
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/postal"
)

type NotifySpace struct {
    courier     postal.CourierInterface
    errorWriter ErrorWriterInterface
}

func NewNotifySpace(courier postal.CourierInterface, errorWriter ErrorWriterInterface) NotifySpace {
    return NotifySpace{
        courier:     courier,
        errorWriter: errorWriter,
    }
}

func (handler NotifySpace) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    params, err := NewNotifyParams(req.Body)
    if err != nil {
        Error(w, 422, []string{"Request body could not be parsed"})
        return
    }

    if !params.Validate() {
        Error(w, 422, params.Errors)
        return
    }

    spaceGUID := strings.TrimPrefix(req.URL.Path, "/spaces/")
    rawToken := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")

    responses, err := handler.courier.Dispatch(rawToken, spaceGUID, postal.IsSpace, params.ToOptions())
    if err != nil {
        handler.errorWriter.Write(w, err)
        return
    }

    output, err := json.Marshal(responses)
    if err != nil {
        panic(err)
    }

    w.WriteHeader(http.StatusOK)
    w.Write(output)
}
