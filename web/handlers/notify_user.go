package handlers

import (
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal/strategies"
	"github.com/cloudfoundry-incubator/notifications/web/params"
	"github.com/ryanmoran/stack"
)

type NotifyUser struct {
	errorWriter ErrorWriterInterface
	notify      NotifyInterface
	strategy    strategies.StrategyInterface
	database    models.DatabaseInterface
}

func NewNotifyUser(notify NotifyInterface, errorWriter ErrorWriterInterface, strategy strategies.StrategyInterface, database models.DatabaseInterface) NotifyUser {
	return NotifyUser{
		errorWriter: errorWriter,
		notify:      notify,
		strategy:    strategy,
		database:    database,
	}
}

func (handler NotifyUser) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	metrics.NewMetric("counter", map[string]interface{}{
		"name": "notifications.web.users",
	}).Log()

	connection := handler.database.Connection()
	err := handler.Execute(w, req, connection, context, handler.strategy)
	if err != nil {
		handler.errorWriter.Write(w, err)
		return
	}
}

func (handler NotifyUser) Execute(w http.ResponseWriter, req *http.Request,
	connection models.ConnectionInterface, context stack.Context, strategy strategies.StrategyInterface) error {
	userGUID := strings.TrimPrefix(req.URL.Path, "/users/")

	output, err := handler.notify.Execute(connection, req, context, userGUID, strategy, params.GUIDValidator{})
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(output)

	return nil
}
