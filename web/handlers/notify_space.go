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
    errorWriter   ErrorWriterInterface
    notify        NotifyInterface
    recipeBuilder RecipeBuilderInterface
    database      models.DatabaseInterface
}

type RecipeBuilderInterface interface {
    NewUAARecipe() postal.UAARecipe
}

func NewNotifySpace(notify NotifyInterface, errorWriter ErrorWriterInterface, recipeBuilder RecipeBuilderInterface, database models.DatabaseInterface) NotifySpace {
    return NotifySpace{
        errorWriter:   errorWriter,
        notify:        notify,
        recipeBuilder: recipeBuilder,
        database:      database,
    }
}

func (handler NotifySpace) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
    connection := handler.database.Connection()
    err := handler.Execute(w, req, connection, context, handler.recipeBuilder.NewUAARecipe())
    if err != nil {
        handler.errorWriter.Write(w, err)
        return
    }

    metrics.NewMetric("counter", map[string]interface{}{
        "name": "notifications.web.spaces",
    }).Log()
}

func (handler NotifySpace) Execute(w http.ResponseWriter, req *http.Request, connection models.ConnectionInterface,
    context stack.Context, recipe postal.UAARecipe) error {

    spaceGUID := postal.SpaceGUID(strings.TrimPrefix(req.URL.Path, "/spaces/"))

    output, err := handler.notify.Execute(connection, req, context, spaceGUID, recipe)
    if err != nil {
        return err
    }

    w.WriteHeader(http.StatusOK)
    w.Write(output)

    return nil
}
