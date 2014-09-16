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
}

type RecipeBuilderInterface interface {
    NewUAARecipe() postal.UAARecipe
}

func NewNotifySpace(notify NotifyInterface, errorWriter ErrorWriterInterface, recipeBuilder RecipeBuilderInterface) NotifySpace {
    return NotifySpace{
        errorWriter:   errorWriter,
        notify:        notify,
        recipeBuilder: recipeBuilder,
    }
}

func (handler NotifySpace) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
    metrics.NewMetric("counter", map[string]interface{}{
        "name": "notifications.web.spaces",
    }).Log()

    connection := models.Database().Connection()
    err := handler.Execute(w, req, connection, handler.recipeBuilder.NewUAARecipe())
    if err != nil {
        handler.errorWriter.Write(w, err)
        return
    }
}

func (handler NotifySpace) Execute(w http.ResponseWriter, req *http.Request,
    connection models.ConnectionInterface, recipe postal.UAARecipe) error {

    spaceGUID := postal.SpaceGUID(strings.TrimPrefix(req.URL.Path, "/spaces/"))

    output, err := handler.notify.Execute(connection, req, spaceGUID, recipe)
    if err != nil {
        return err
    }

    w.WriteHeader(http.StatusOK)
    w.Write(output)

    return nil
}
