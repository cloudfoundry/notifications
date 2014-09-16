package handlers

import (
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/metrics"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/ryanmoran/stack"
)

type NotifyUser struct {
    errorWriter   ErrorWriterInterface
    notify        NotifyInterface
    recipeBuilder RecipeBuilderInterface
}

func NewNotifyUser(notify NotifyInterface, errorWriter ErrorWriterInterface, recipeBuilder RecipeBuilderInterface) NotifyUser {
    return NotifyUser{
        errorWriter:   errorWriter,
        notify:        notify,
        recipeBuilder: recipeBuilder,
    }
}

func (handler NotifyUser) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
    connection := models.Database().Connection()
    err := handler.Execute(w, req, connection, context, handler.recipeBuilder.NewUAARecipe())
    if err != nil {
        handler.errorWriter.Write(w, err)
        return
    }

    metrics.NewMetric("counter", map[string]interface{}{
        "name": "notifications.web.users",
    }).Log()
}

func (handler NotifyUser) Execute(w http.ResponseWriter, req *http.Request,
    connection models.ConnectionInterface, context stack.Context, recipe postal.UAARecipe) error {
    userGUID := postal.UserGUID(strings.TrimPrefix(req.URL.Path, "/users/"))

    output, err := handler.notify.Execute(connection, req, context, userGUID, recipe)
    if err != nil {
        return err
    }

    w.WriteHeader(http.StatusOK)
    w.Write(output)

    return nil
}
