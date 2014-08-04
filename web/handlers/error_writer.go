package handlers

import (
    "encoding/json"
    "net/http"

    "github.com/cloudfoundry-incubator/notifications/postal"
)

type ErrorWriterInterface interface {
    Write(http.ResponseWriter, error)
}

type ErrorWriter struct{}

func NewErrorWriter() ErrorWriter {
    return ErrorWriter{}
}

func (handler ErrorWriter) Write(w http.ResponseWriter, err error) {
    switch err.(type) {
    case postal.CCDownError:
        Error(w, http.StatusBadGateway, []string{"Cloud Controller is unavailable"})
    case postal.CCNotFoundError:
        Error(w, http.StatusNotFound, []string{err.Error()})
    case postal.UAADownError:
        Error(w, http.StatusBadGateway, []string{"UAA is unavailable"})
    case postal.UAAGenericError:
        Error(w, http.StatusBadGateway, []string{err.Error()})
    case postal.TemplateLoadError:
        Error(w, http.StatusInternalServerError, []string{"An email template could not be loaded"})
    default:
        panic(err)
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
