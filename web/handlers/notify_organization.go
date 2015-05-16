package handlers

import (
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal/strategies"
	"github.com/cloudfoundry-incubator/notifications/web/params"
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
	err := handler.Execute(w, req, connection, context)
	if err != nil {
		handler.errorWriter.Write(w, err)
		return
	}
}

func (handler NotifyOrganization) Execute(w http.ResponseWriter, req *http.Request, connection models.ConnectionInterface,
	context stack.Context) error {

	organizationGUID := strings.TrimPrefix(req.URL.Path, "/organizations/")

	vcapRequestID := context.Get(VCAPRequestIDKey).(string)

	output, err := handler.notify.Execute(connection, req, context, organizationGUID, handler.strategy, params.GUIDValidator{}, vcapRequestID)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(output)

	return nil
}
