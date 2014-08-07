package handlers

import (
    "encoding/json"
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/postal"
)

type NotifyUser struct {
    courier     postal.CourierInterface
    errorWriter ErrorWriterInterface
}

func NewNotifyUser(courier postal.CourierInterface, errorWriter ErrorWriterInterface) NotifyUser {
    return NotifyUser{
        courier:     courier,
        errorWriter: errorWriter,
    }
}

func (handler NotifyUser) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    userGUID := strings.TrimPrefix(req.URL.Path, "/users/")

    params, err := NewNotifyParams(req.Body)
    if err != nil {
        handler.errorWriter.Write(w, err)
        return
    }

    if !params.Validate() {
        handler.errorWriter.Write(w, ParamsValidationError(params.Errors))
        return
    }

    rawToken := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")

    responses, err := handler.courier.Dispatch(rawToken, postal.UserGUID(userGUID), params.ToOptions())
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
