package handlers

import (
    "encoding/json"
    "net/http"
    "strings"

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

    notificationType := notificationType(templateName)

    if notificationType != "" {
        template, err := handler.Finder.Find(notificationType, templateName)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        response, err := json.Marshal(template)

        w.Write(response)
    } else {
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte("Could not find template. Did you specify a notification type?"))
    }
}

func notificationType(URL string) string {
    if strings.HasSuffix(URL, services.UserBody) {
        return services.UserBody
    } else if strings.HasSuffix(URL, services.SpaceBody) {
        return services.SpaceBody
    }
    return ""
}
