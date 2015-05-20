package handlers

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal/strategies"
	"github.com/cloudfoundry-incubator/notifications/web/params"
	"github.com/ryanmoran/stack"
)

type NotifyEveryone struct {
	errorWriter ErrorWriterInterface
	notify      NotifyInterface
	strategy    strategies.StrategyInterface
}

func NewNotifyEveryone(notify NotifyInterface, errorWriter ErrorWriterInterface, strategy strategies.StrategyInterface) NotifyEveryone {
	return NotifyEveryone{
		errorWriter: errorWriter,
		notify:      notify,
		strategy:    strategy,
	}
}

func (handler NotifyEveryone) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	connection := context.Get("database").(models.DatabaseInterface).Connection()
	vcapRequestID := context.Get(VCAPRequestIDKey).(string)

	output, err := handler.notify.Execute(connection, req, context, "", handler.strategy, params.GUIDValidator{}, vcapRequestID)
	if err != nil {
		handler.errorWriter.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(output)
}
