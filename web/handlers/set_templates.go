package handlers

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/web/services"
    "github.com/ryanmoran/stack"
)

type SetTemplates struct {
    updater     services.TemplateUpdaterInterface
    ErrorWriter ErrorWriterInterface
}

func NewSetTemplates(updater services.TemplateUpdaterInterface, errorWriter ErrorWriterInterface) SetTemplates {
    return SetTemplates{
        updater:     updater,
        ErrorWriter: errorWriter,
    }
}

func (handler SetTemplates) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
    var template models.Template

    templateName := strings.Split(req.URL.String(), "/templates/")[1]

    respBody, err := ioutil.ReadAll(req.Body)
    if err != nil {
        panic(err)
    }

    valid, err := handler.requestIsValid(respBody)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    if !valid {
        w.WriteHeader(422)
        return
    }

    err = json.Unmarshal(respBody, &template)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    template.Name = templateName
    template.Overridden = true
    err = handler.updater.Update(template)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func (handler SetTemplates) requestIsValid(body []byte) (bool, error) {
    var templateMap map[string]string

    err := json.Unmarshal(body, &templateMap)
    if err != nil {
        return false, err
    }

    _, textExists := templateMap["text"]
    _, htmlExists := templateMap["html"]

    return (textExists && htmlExists), nil
}
