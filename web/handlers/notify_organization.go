package handlers

import (
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/metrics"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/cloudfoundry-incubator/notifications/postal/strategies"
    "github.com/ryanmoran/stack"
)

type NotifyOrganization struct {
    errorWriter ErrorWriterInterface
    notify      NotifyInterface
    strategy    strategies.StrategyInterface
    database    models.DatabaseInterface
}

func NewNotifyOrganization(notify NotifyInterface, errorWriter ErrorWriterInterface, strategy strategies.StrategyInterface, database models.DatabaseInterface) NotifyOrganization {
    return NotifyOrganization{
        errorWriter: errorWriter,
        notify:      notify,
        strategy:    strategy,
        database:    database,
    }
}

func (handler NotifyOrganization) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
    connection := handler.database.Connection()
    err := handler.Execute(w, req, connection, context, handler.strategy)
    if err != nil {
        handler.errorWriter.Write(w, err)
        return
    }

    metrics.NewMetric("counter", map[string]interface{}{
        "name": "notifications.web.organizations",
    }).Log()
}

func (handler NotifyOrganization) Execute(w http.ResponseWriter, req *http.Request, connection models.ConnectionInterface,
    context stack.Context, strategy strategies.StrategyInterface) error {

    organizationGUID := postal.UAAGUID(strings.TrimPrefix(req.URL.Path, "/organizations/"))

    output, err := handler.notify.Execute(connection, req, context, organizationGUID, strategy)
    if err != nil {
        return err
    }

    w.WriteHeader(http.StatusOK)
    w.Write(output)

    return nil
}
