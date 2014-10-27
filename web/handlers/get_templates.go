package handlers

import (
    "encoding/json"
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/web/services"
    "github.com/ryanmoran/stack"
)

type GetTemplates struct {
    Finder      services.TemplateFinderInterface
    ErrorWriter ErrorWriterInterface
}

func NewGetTemplates(templateFinder services.TemplateFinderInterface, errorWriter ErrorWriterInterface) GetTemplates {
    return GetTemplates{
        Finder:      templateFinder,
        ErrorWriter: errorWriter,
    }
}

func (handler GetTemplates) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
    templateName := strings.Split(req.URL.Path, "/templates/")[1]

    template, err := handler.Finder.Find(templateName)
    if err != nil {
        handler.ErrorWriter.Write(w, err)
        return
    }

    response, err := json.Marshal(template)
    if err != nil {
        panic(err)
    }
    w.Write(response)
}
