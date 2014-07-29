package handlers

import (
    "encoding/json"
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/postal"
)

type NotifyUser struct {
    courier postal.CourierInterface
}

func NewNotifyUser(courier postal.CourierInterface) NotifyUser {
    return NotifyUser{
        courier: courier,
    }
}

func (handler NotifyUser) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    userGUID := strings.TrimPrefix(req.URL.Path, "/users/")

    params, err := NewNotifyParams(req.Body)
    if err != nil {
        Error(w, 422, []string{"Request body could not be parsed"})
        return
    }

    if !params.Validate() {
        Error(w, 422, params.Errors)
        return
    }

    rawToken := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")

    responses, err := handler.courier.Dispatch(rawToken, userGUID, postal.IsUser, params.ToOptions())
    if err != nil {
        switch err.(type) {
        case postal.CCDownError:
            Error(w, http.StatusBadGateway, []string{"Cloud Controller is unavailable"})
        case postal.UAADownError:
            Error(w, http.StatusBadGateway, []string{"UAA is unavailable"})
        case postal.UAAGenericError:
            Error(w, http.StatusBadGateway, []string{err.Error()})
        case postal.TemplateLoadError:
            Error(w, http.StatusInternalServerError, []string{"An email template could not be loaded"})
        default:
            panic(err)
        }
        return
    }

    output, err := json.Marshal(responses)
    if err != nil {
        panic(err)
    }

    w.WriteHeader(http.StatusOK)
    w.Write(output)
}
