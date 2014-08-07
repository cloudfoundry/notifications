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

func (writer ErrorWriter) Write(w http.ResponseWriter, err error) {
    switch err.(type) {
    case postal.CCDownError:
        writer.write(w, http.StatusBadGateway, []string{"Cloud Controller is unavailable"})
    case postal.CCNotFoundError:
        writer.write(w, http.StatusNotFound, []string{err.Error()})
    case postal.UAADownError:
        writer.write(w, http.StatusBadGateway, []string{"UAA is unavailable"})
    case postal.UAAGenericError:
        writer.write(w, http.StatusBadGateway, []string{err.Error()})
    case postal.TemplateLoadError:
        writer.write(w, http.StatusInternalServerError, []string{"An email template could not be loaded"})
    case ParamsParseError:
        writer.write(w, 422, []string{err.Error()})
    case ParamsValidationError:
        writer.write(w, 422, err.(ParamsValidationError).Errors())
    default:
        panic(err)
    }
}

func (writer ErrorWriter) write(w http.ResponseWriter, code int, errors []string) {
    response, err := json.Marshal(map[string][]string{
        "errors": errors,
    })
    if err != nil {
        panic(err)
    }

    w.WriteHeader(code)
    w.Write(response)
}
