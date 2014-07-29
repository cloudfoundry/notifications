package handlers

import (
    "encoding/json"
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/postal"
)

type NotifySpace struct {
    courier postal.CourierInterface
}

func NewNotifySpace(courier postal.CourierInterface) NotifySpace {
    return NotifySpace{
        courier: courier,
    }
}

func Error(w http.ResponseWriter, code int, errors []string) {
    response, err := json.Marshal(map[string][]string{
        "errors": errors,
    })
    if err != nil {
        panic(err)
    }

    w.WriteHeader(code)
    w.Write(response)
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
