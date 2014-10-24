package handlers

import (
    "encoding/json"
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/web/services"
    "github.com/ryanmoran/stack"
)

type GetTemplates struct {
    Finder services.TemplateFinderInterface
}

func NewGetTemplates(templateFinder services.TemplateFinderInterface) GetTemplates {
    return GetTemplates{Finder: templateFinder}
}

func (handler GetTemplates) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
    templateName := strings.Split(req.URL.Path, "/templates/")[1]

    template, err := handler.Finder.Find(templateName)
    if err != nil {
        if (err == models.ErrRecordNotFound{}) {
            w.WriteHeader(http.StatusNotFound)
            w.Write([]byte("Could not find template. Did you specify a notification type?"))
        } else {
            w.WriteHeader(http.StatusInternalServerError)
        }
        return
    }

    response, err := json.Marshal(template)
    if err != nil {
        panic(err)
    }
    w.Write(response)
}
