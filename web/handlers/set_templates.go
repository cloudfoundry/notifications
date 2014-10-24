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
    updater services.TemplateUpdaterInterface
}

func NewSetTemplates(updater services.TemplateUpdaterInterface) SetTemplates {
    return SetTemplates{
        updater: updater,
    }
}

func (handler SetTemplates) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
    var template models.Template
    templateName := strings.Split(req.URL.String(), "/templates/")[1]

    respBody, err := ioutil.ReadAll(req.Body)
    if err != nil {
        panic(err)
    }

    err = json.Unmarshal(respBody, &template)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    if template.Text == "" || template.HTML == "" {
        w.WriteHeader(422)
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
