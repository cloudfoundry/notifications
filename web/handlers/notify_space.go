package handlers

import (
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/metrics"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/ryanmoran/stack"
)

type NotifySpace struct {
    errorWriter ErrorWriterInterface
    notify      NotifyInterface
    recipe      postal.RecipeInterface
    database    models.DatabaseInterface
}

func NewNotifySpace(notify NotifyInterface, errorWriter ErrorWriterInterface, recipe postal.RecipeInterface, database models.DatabaseInterface) NotifySpace {
    return NotifySpace{
        errorWriter: errorWriter,
        notify:      notify,
        recipe:      recipe,
        database:    database,
    }
}

func (handler NotifySpace) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
    connection := handler.database.Connection()
    err := handler.Execute(w, req, connection, context, handler.recipe)
    if err != nil {
        handler.errorWriter.Write(w, err)
        return
    }

    metrics.NewMetric("counter", map[string]interface{}{
        "name": "notifications.web.spaces",
    }).Log()
}

func (handler NotifySpace) Execute(w http.ResponseWriter, req *http.Request, connection models.ConnectionInterface,
    context stack.Context, recipe postal.RecipeInterface) error {

    spaceGUID := postal.SpaceGUID(strings.TrimPrefix(req.URL.Path, "/spaces/"))

    output, err := handler.notify.Execute(connection, req, context, spaceGUID, recipe)
    if err != nil {
        return err
    }

    w.WriteHeader(http.StatusOK)
    w.Write(output)

    return nil
}
